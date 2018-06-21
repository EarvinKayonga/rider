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
	"github.com/EarvinKayonga/rider/httpx"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/messaging"
	"github.com/EarvinKayonga/rider/stats"
)

// NewGatewayService returns the gateway service wrapped in a valid http.Server.
func NewGatewayService(
	ctx context.Context,
	m Metadata,
	conf configuration.GatewayConfiguration,
	logger logging.Logger,
	statter stats.Statter,
	messenger messaging.Emitter) (*http.Server, error) {

	router := mux.NewRouter()
	router.StrictSlash(true)

	limiter := httpx.Limiter(conf.Limiter.Limit, conf.Limiter.Burst)
	err := registerRoutesForGatewayService(ctx, router, m, conf, logger, statter, messenger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register routes for gateway service")
	}

	limitedRouter := limiter(router)
	securedRouter := secureHeaders(limitedRouter)

	return NewServer(ctx, conf.Server, securedRouter)
}

func registerRoutesForGatewayService(
	ctx context.Context,
	router *mux.Router,
	metadata Metadata,
	conf configuration.GatewayConfiguration,
	logger logging.Logger,
	statter stats.Statter,
	messenger messaging.Emitter) error {

	router.HandleFunc("/health", health(ctx, metadata))

	router.HandleFunc("/bike/{bikeID}", GatewayGetBikeByID(ctx, conf, logger, statter))
	router.HandleFunc("/bikes", GatewayListOfBikes(ctx, conf, logger, statter))
	router.HandleFunc("/trip/track", TrackTrip(ctx, conf, logger, statter, messenger))
	router.HandleFunc("/trip/start", GatewayStartTrip(ctx, conf, logger, statter))
	router.HandleFunc("/trip/end", GatewayEndTrip(ctx, conf, logger, statter))

	return nil
}

// GatewayGetBikeByID returns a bike given an ID.
func GatewayGetBikeByID(ctx context.Context,
	conf configuration.GatewayConfiguration,
	logger logging.Logger,
	statter stats.Statter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("gateway.bike.timing", time.Since(start), 1.0)
			_ = statter.Inc("gateway.bike.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bike, err := domain.GatewayGetBikeByID(ctx, conf, mux.Vars(req)["bikeID"])
		if err != nil {
			defer func() {
				_ = statter.Inc("gateway.bike.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while fetching bike from bike service")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bike)
		if err != nil {
			defer func() {
				_ = statter.Inc("gateway.bike.error", 1, 1.0)
			}()

			Erroring(ctx, w, err, logger)
			logger.WithError(err).Error("an error occuring while rendering bike")
		}

		logger.Info("bike successfully rendered")
	}
}

// GatewayListOfBikes returns a paginated list of bikes.
func GatewayListOfBikes(ctx context.Context,
	conf configuration.GatewayConfiguration,
	logger logging.Logger,
	statter stats.Statter,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("gateway.list.bikes.timing", time.Since(start), 1.0)
			_ = statter.Inc("gateway.list.bikes.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cursor, limit := GetPaginationArguments(req)
		bikes, err := domain.GatewayListOfBikes(ctx, conf, cursor, limit)
		if err != nil {
			defer func() {
				_ = statter.Inc("gateway.list.bikes.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while fetching list of bikes")
			Erroring(ctx, w, err, logger)
			return
		}

		err = json.NewEncoder(w).Encode(bikes)
		if err != nil {
			defer func() {
				_ = statter.Inc("gateway.list.bikes.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while rendering list of bikes")
			Erroring(ctx, w, err, logger)
			return
		}

		logger.Info("bike list successfully rendered")
	}
}

// GatewayStartTrip is the handler for starting a trip.
func GatewayStartTrip(
	ctx context.Context,
	conf configuration.GatewayConfiguration,
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

		trip, err := domain.GatewayStartTrip(ctx, conf, tripPayload.BikeID, tripPayload.Location.Lat,
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

// GatewayEndTrip is the handler for starting a trip.
func GatewayEndTrip(
	ctx context.Context,
	conf configuration.GatewayConfiguration,
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

		trip, err := domain.GatewayEndTrip(ctx, conf, tripPayload.TripID, tripPayload.Location.Lat,
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

// TrackTrip is the middleware for tracking a trip.
// By tracking, we mean here, adding a location to a trip.
func TrackTrip(ctx context.Context,
	conf configuration.GatewayConfiguration,
	logger logging.Logger,
	statter stats.Statter,
	messenger messaging.Emitter) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			_ = statter.TimingDuration("track.trip.timing", time.Since(start), 1.0)
			_ = statter.Inc("track.trip.request", 1, 1.0)
		}()

		if req == nil {
			logger.Error("received an empty request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer func() {
			_ = req.Body.Close()
		}()

		hearbeat := domain.TrackTripPayload{}

		err := json.NewDecoder(req.Body).Decode(&hearbeat)
		if err != nil {
			defer func() {
				_ = statter.Inc("track.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while unmarshalling end payload")
			Erroring(ctx, w, err, logger)
			return
		}

		err = domain.TrackTrip(ctx, messenger, hearbeat)
		if err != nil {
			defer func() {
				_ = statter.Inc("track.trip.error", 1, 1.0)
			}()

			logger.WithError(err).Error("an error occuring while sending track trip")
			Erroring(ctx, w, err, logger)
			return
		}

		logger.Info("location successfully sent")
	}

}
