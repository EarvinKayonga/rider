package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/domain"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/stats"
)

// NewBikeService returns the bike service wrapped in a valid http.Server.
func NewBikeService(
	ctx context.Context,
	m Metadata,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter) (*http.Server, error) {

	router := mux.NewRouter()
	router.StrictSlash(true)

	err := registerRoutesForBikeService(ctx, router, m, conf, logger, statter)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register routes for bike service")
	}

	securedRouter := secureHeaders(router)
	return NewServer(ctx, conf.Server, securedRouter)
}

func registerRoutesForBikeService(
	ctx context.Context,
	router *mux.Router,
	metadata Metadata,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter) error {

	router.HandleFunc("/health", health(ctx, metadata)).Methods("GET")
	router.HandleFunc("/bike/{bikeID}", GetBikeByID(ctx, conf, logger, statter)).Methods("GET")
	router.HandleFunc("/lock/{bikeID}", LockBikeByID(ctx, conf, logger, statter)).Methods("GET")
	router.HandleFunc("/unlock/{bikeID}", UnLockBikeByID(ctx, conf, logger, statter)).Methods("GET")
	router.HandleFunc("/bikes", ListOfBikes(ctx, conf, logger, statter)).Methods("GET")

	return nil
}

// ListOfBikes returns a paginated list of bikes.
func ListOfBikes(ctx context.Context,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("list.bikes.timing", time.Since(start), 1.0)
			_ = statter.Inc("list.bikes.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cursor, limit := GetPaginationArguments(req)
		bikes, err := domain.ListOfBikes(ctx, cursor, limit)
		if err != nil {
			defer func() {
				_ = statter.Inc("list.bikes.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while fetching list of bikes")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bikes)
		if err != nil {
			defer func() {
				_ = statter.Inc("list.bikes.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while rendering list of bikes")
			Erroring(ctx, w, err, logger)
			return
		}

		logger.Info("bike list successfully rendered")
	}
}

// GetBikeByID returns a bike given an ID.
func GetBikeByID(ctx context.Context,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("bike.timing", time.Since(start), 1.0)
			_ = statter.Inc("bike.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params := mux.Vars(req)
		bikeID := params["bikeID"]

		bike, err := domain.GetBikeByID(ctx, bikeID)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while fetching bike")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bike)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering bike")
		}

		logger.Info("bike successfully rendered")
	}
}

// LockBikeByID lock a bike given an ID.
func LockBikeByID(ctx context.Context,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("bike.lock.timing", time.Since(start), 1.0)
			_ = statter.Inc("bike.lock.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params := mux.Vars(req)
		bikeID := params["bikeID"]

		bike, err := domain.LockBikeByID(ctx, bikeID)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.lock.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while locking bike")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bike)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.lock.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering bike")
		}

		logger.Info("locked bike successfully rendered")
	}
}

// UnLockBikeByID lock a bike given an ID.
func UnLockBikeByID(ctx context.Context,
	conf configuration.BikeConfiguration,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("bike.unlock.timing", time.Since(start), 1.0)
			_ = statter.Inc("bike.unlock.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params := mux.Vars(req)
		bikeID := params["bikeID"]

		bike, err := domain.UnLockBikeByID(ctx, bikeID)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.unlock.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while unlocking bike")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bike)
		if err != nil {
			defer func() {
				_ = statter.Inc("bike.unlock.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering bike")
		}

		logger.Info("unlocked bike successfully rendered")
	}
}
