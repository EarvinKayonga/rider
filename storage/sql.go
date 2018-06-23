package storage

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/models"
)

// scannable for scanning rows.
type scannable interface {
	Scan(...interface{}) error
}

// Bike Queries.
var (
	findBikeByPublicID = `SELECT id, public_id, latitude, longitude, status FROM bikes WHERE public_id=$1;`
	listBikes          = `SELECT id, public_id, latitude, longitude, status FROM bikes WHERE public_id <= $2 ORDER BY public_id DESC LIMIT $1;`
	listAllBikes       = `SELECT id, public_id, latitude, longitude, status FROM bikes LIMIT $1;`
	updateBikeLocation = `UPDATE bikes SET latitude = $2, longitude = $3 WHERE public_id = $1;`
	createBike         = `INSERT INTO bikes (public_id, latitude, longitude, status) VALUES ($1, $2, $3, $4) 
							RETURNING id, public_id, latitude, longitude, status;`

	unlockBikeByPublicID = `UPDATE bikes SET status = 1 
							WHERE public_id = $1
							RETURNING id, public_id, latitude, longitude, status;`

	lockBikeByPublicID = `UPDATE bikes SET status = 0
							WHERE public_id = $1
							RETURNING id, public_id, latitude, longitude, status;`
)

// toBike centralizes the parsing of a sql Row to a models.Bike.
func toBike(row scannable) (*models.Bike, error) {
	var publicID string
	var id int64
	var status int
	var latitude, longitude float64

	err := row.Scan(&id, &publicID, &latitude, &longitude, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBikeNotFound
		}

		return nil, errors.Wrap(err,
			"an error occured while scanning for bike")
	}

	return &models.Bike{
		ID:       publicID,
		Status:   status,
		Location: models.CreateLocation(latitude, longitude),
	}, nil
}

// Location queries.
var (
	addLocationToTrip = `INSERT INTO locations (latitude, longitude, trip_id) VALUES ($1, $2, $3) 
							RETURNING id, latitude, longitude, trip_id, created_at;`

	listLocationForTrip = `SELECT id, latitude, longitude, trip_id, created_at FROM locations WHERE trip_id=$1;`
)

// toLocation centralizes the parsing of a sql Row to a Location.
func toLocation(row scannable) (*Location, error) {
	var tripID string
	var createdAt time.Time
	var id int64
	var latitude, longitude float64

	err := row.Scan(&id, &latitude, &longitude, &tripID, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBikeNotFound
		}

		return nil, errors.Wrap(err,
			"an error occured while scanning for location")
	}

	return &Location{
		ID:        id,
		Longitude: longitude,
		Latitude:  latitude,
		CreatedAt: createdAt,
		TripID:    tripID,
	}, nil
}

// Trip queries
var (
	createTrip = `INSERT INTO trips (bike_id, public_id , status) VALUES($1, $2, $3)
					RETURNING id, started_at, ended_at, public_id, bike_id, status;`

	endTrip = `UPDATE trips SET status = 0,  public_id = $1
					RETURNING id, started_at, ended_at, public_id, bike_id, status;`

	listTrip = `SELECT id, started_at, ended_at, public_id, bike_id, status FROM trips WHERE public_id=$1;`
)

// toTrip centralizes the parsing of a sql Row to a Trip.
func toTrip(row scannable) (*Trip, error) {
	var bikeID, publicID string
	var startedAt time.Time
	var endedAt pq.NullTime
	var id int64
	var status int

	err := row.Scan(&id, &startedAt, &endedAt, &publicID, &bikeID, &status)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while scanning for location")
	}

	return &Trip{
		ID:        id,
		BikeID:    bikeID,
		StartedAt: startedAt,
		PublicID:  publicID,
		EndedAt:   endedAt,
		Status:    status,
	}, nil
}

// unwrapNullTime converts a null time to a time.
func unwrapNullTime(t pq.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}

	return nil
}
