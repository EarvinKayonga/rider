package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	// loading postgres drivers
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/entropy"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/models"
)

// pgStore implements Store
// with postgres.
type pgStore struct {
	database *sql.DB
	logger   logging.Logger
}

// NewPostgresDatabase return a valid Store
// with postgres.
func NewPostgresDatabase(ctx context.Context,
	config configuration.Database, logger logging.Logger) (Store, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port,
		config.User, config.Password, config.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while connecting to database with %v",
			config)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while pinging database")
	}

	logger.Info("succesfully opened an connection and pinged database")

	return &pgStore{
		database: db,
		logger:   logger,
	}, nil
}

func (e *pgStore) CreateBikes(ctx context.Context, bikes []Bike) ([]models.Bike, error) {
	created := []models.Bike{}
	for _, b := range bikes {
		bike, err := toBike(e.
			database.
			QueryRow(createBike, b.PublicID, b.Latitude, b.Longitude, b.Status))

		if err != nil {
			return nil, errors.Wrap(err, "an error occured while inserting a bike in db")
		}

		created = append(created, *bike)
	}

	e.logger.Infof("inserted %d bikes", len(created))

	return created, nil
}

func (e *pgStore) UpdateBikeLocation(ctx context.Context, bikeID string, lat, lng float64) error {
	result, err := e.database.Exec(updateBikeLocation, bikeID, lat, lng)
	if err != nil {
		errors.Wrap(err,
			"an error occured while updating bike location")
	}

	count, err := result.RowsAffected()
	if err != nil {
		errors.Wrap(err,
			"an error occured while checking nbs of affected rows through updating bike location")
	}

	if count == 0 {
		return ErrBikeNotFound
	}

	e.logger.Infof("location of %d bikes updated", count)
	return nil
}

func (e *pgStore) FindBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return toBike(e.
		database.
		QueryRow(findBikeByPublicID, bikeID))
}

func (e *pgStore) UnLockBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return toBike(e.
		database.
		QueryRow(unlockBikeByPublicID, bikeID))
}

func (e *pgStore) LockBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return toBike(e.
		database.
		QueryRow(lockBikeByPublicID, bikeID))
}

func (e *pgStore) ListAllBikes(ctx context.Context, limit int64) ([]models.Bike, error) {
	rows, err := e.database.Query(listAllBikes, limit)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while listing bikes from the database")
	}

	if rows != nil {
		defer func() {
			thr := rows.Close()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while closing query")
			}
		}()
	}

	bikes := []models.Bike{}

	for rows.Next() {
		bike, err := toBike(rows)
		if err != nil {
			return nil, errors.Wrap(err,
				"an error occured while querying a bike")
		}

		bikes = append(bikes, *bike)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured during iteration over returned bikes")
	}

	e.logger.Infof("returned %d bikes", len(bikes))

	return bikes, nil
}

func (e *pgStore) optionsList(ctx context.Context, cursor string, limit int64) (*sql.Rows, error) {

	if cursor == "" {
		return e.database.Query(listAllBikes, limit)
	}

	if limit == 0 {
		limit = 20
	}

	return e.database.Query(listBikes, limit, cursor)
}

func (e *pgStore) ListBikes(ctx context.Context, cursor string, limit int64) ([]models.Bike, error) {
	rows, err := e.optionsList(ctx, cursor, limit)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while listing bikes from the database")
	}

	if rows != nil {
		defer func() {
			thr := rows.Close()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while closing query")
			}
		}()
	}

	bikes := []models.Bike{}

	for rows.Next() {
		bike, err := toBike(rows)
		if err != nil {
			return nil, errors.Wrap(err,
				"an error occured while querying a bike")
		}

		bikes = append(bikes, *bike)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured during iteration over returned bikes")
	}

	e.logger.Infof("returned %d bikes", len(bikes))

	return bikes, nil
}

func (e *pgStore) AddLocationToTrip(ctx context.Context, tripID string, lat, lng float64) error {
	_, err := toLocation(
		e.
			database.
			QueryRow(addLocationToTrip, lat, lng, tripID))

	if err != nil {
		return errors.Wrap(err,
			"an error occured while adding location to trip")
	}

	e.logger.Info("added location to trip")

	return nil
}

func (e *pgStore) CreateTrip(ctx context.Context, bikeID string, lat, lng float64) (*models.Trip, error) {
	tx, err := e.database.Begin()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while begin transaction for trip creation")
	}

	stmt, err := tx.Prepare(createTrip)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while preparing transaction for trip creation")
	}

	if stmt != nil {
		defer func() {
			thr := stmt.Close()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while closing prepared statement")
			}
		}()
	}

	trip, err := toTrip(stmt.QueryRow(bikeID, entropy.FromContext(ctx).NewID(), 1))
	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing trip in database")
	}

	bike, err := toBike(
		e.
			database.
			QueryRow(lockBikeByPublicID, bikeID))

	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing locking bike in database")
	}

	location, err := toLocation(
		e.
			database.
			QueryRow(addLocationToTrip, lat, lng, trip.PublicID))

	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing location in database")
	}

	err = tx.Commit()
	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while committing a transaction")
	}

	e.logger.Info("successfully created trip")

	return &models.Trip{
		Locations: []models.Location{
			models.CreateLocation(location.Latitude, location.Longitude),
		},

		ID:        trip.PublicID,
		Status:    trip.Status,
		BikeID:    bike.ID,
		StartedAt: trip.StartedAt,
		EndedAt:   unwrapNullTime(trip.EndedAt),
	}, nil
}

func (e *pgStore) EndTrip(ctx context.Context, tripID string, lat, lng float64) (*models.Trip, error) {
	tx, err := e.database.Begin()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while begin transaction for trip ending")
	}

	stmt, err := tx.Prepare(endTrip)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while preparing transaction for trip ending")
	}

	if stmt != nil {
		defer func() {
			thr := stmt.Close()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while closing prepared statement")
			}
		}()
	}

	trip, err := toTrip(stmt.QueryRow(tripID))
	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing trip in database")
	}

	_, err = toLocation(
		e.
			database.
			QueryRow(addLocationToTrip, lat, lng, trip.PublicID))

	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing location in database")
	}

	bike, err := toBike(
		e.
			database.
			QueryRow(unlockBikeByPublicID, trip.BikeID))

	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while writing unlocking bike in database")
	}

	err = tx.Commit()
	if err != nil {
		defer func() {
			thr := tx.Rollback()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while rollbacking transaction")
			}
		}()

		return nil, errors.Wrap(err,
			"an error occured while committing a transaction")
	}

	locations, err := e.GetLocationsForTrip(ctx, trip.PublicID)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while fetching locations for trip: %s", trip.ID)
	}

	e.logger.Info("successfully ended trip")

	return &models.Trip{
		Locations: locations,

		ID:        trip.PublicID,
		Status:    trip.Status,
		BikeID:    bike.ID,
		StartedAt: trip.StartedAt,
		EndedAt:   unwrapNullTime(trip.EndedAt),
	}, nil
}

func (e *pgStore) GetLocationsForTrip(ctx context.Context, tripID string) ([]models.Location, error) {
	rows, err := e.database.Query(listTrip, tripID)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while listing trips from the database")
	}

	if rows != nil {
		defer func() {
			thr := rows.Close()
			if thr != nil {
				e.logger.WithError(thr).Warn("error while closing location list query")
			}
		}()
	}

	locations := []Location{}
	for rows.Next() {
		location, err := toLocation(rows)
		if err != nil {
			return nil, errors.Wrap(err,
				"an error occured while querying a location")
		}

		locations = append(locations, *location)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured during iteration over returned locations")
	}

	count := len(locations)

	correctLocations := make([]models.Location, count)
	e.logger.Infof("fetched %s location points", count)

	sort.SliceStable(locations, func(i, j int) bool {
		return locations[i].CreatedAt.Unix() < locations[j].CreatedAt.Unix()
	})

	for _, location := range locations {
		correctLocations = append(correctLocations,
			models.CreateLocation(location.Latitude, location.Longitude))
	}

	return correctLocations, nil
}

func (e *pgStore) Close(_ context.Context) error {
	return e.database.Close()
}
