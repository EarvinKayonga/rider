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

// NewTripService returns the trip service wrapped in a valid http.Server.
func NewTripService(
	ctx context.Context,
	m Metadata,
	conf configuration.TripConfiguration,
	logger logging.Logger,
	statter stats.Statter) (*http.Server, error) {

	router := mux.NewRouter()
	router.StrictSlash(true)

	err := registerRoutesForTripService(ctx, router, m, conf, logger, statter)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register routes for Trip service")
	}

	securedRouter := secureHeaders(router)

	return NewServer(ctx, conf.Server, securedRouter)
}

func registerRoutesForTripService(
	ctx context.Context,
	router *mux.Router,
	metadata Metadata,
	conf configuration.TripConfiguration,
	logger logging.Logger,
	statter stats.Statter) error {

	router.HandleFunc("/health", health(ctx, metadata)).Methods("GET")
	router.HandleFunc("/trip/start", StartTrip(ctx, logger, statter)).Methods("POST", "PUT")
	router.HandleFunc("/trip/end", EndTrip(ctx, logger, statter)).Methods("POST", "PUT")

	return nil
}

// StartTrip is the handler for starting a trip.
func StartTrip(
	ctx context.Context,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("start.trip.timing", time.Since(start), 1.0)
			_ = statter.Inc("start.trip.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() {
			_ = req.Body.Close()
		}()

		tripPayload := domain.StartTripPayload{}

		err := json.NewDecoder(req.Body).Decode(&tripPayload)
		if err != nil {
			defer func() {
				_ = statter.Inc("start.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while unmarshalling start payload")
			Erroring(ctx, w, err, logger)
			return
		}

		trip, err := domain.StartTrip(ctx, tripPayload.BikeID, tripPayload.Location.Lat,
			tripPayload.Location.Lng)
		if err != nil {
			defer func() {
				_ = statter.Inc("start.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while starting trip")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(trip)
		if err != nil {
			defer func() {
				_ = statter.Inc("start.trip.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering trip")
			return
		}

		logger.Info("bike successfully started")
	}
}

// EndTrip is the handler for starting a trip.
func EndTrip(
	ctx context.Context,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("end.trip.timing", time.Since(start), 1.0)
			_ = statter.Inc("end.trip.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() {
			_ = req.Body.Close()
		}()

		tripPayload := domain.EndTripPayload{}

		err := json.NewDecoder(req.Body).Decode(&tripPayload)
		if err != nil {
			defer func() {
				_ = statter.Inc("end.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while unmarshalling end payload")
			Erroring(ctx, w, err, logger)
			return
		}

		trip, err := domain.EndTrip(ctx, tripPayload.TripID, tripPayload.Location.Lat,
			tripPayload.Location.Lng)
		if err != nil {
			defer func() {
				_ = statter.Inc("end.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while ending trip")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(trip)
		if err != nil {
			defer func() {
				_ = statter.Inc("end.trip.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering trip")
			return
		}

		logger.Info("trip successfully ended")
	}
}
