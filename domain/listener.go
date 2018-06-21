package domain

import (
	"context"

	"github.com/pkg/errors"
	bus "github.com/rafaeljesus/nsq-event-bus"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/logging"
	"github.com/EarvinKayonga/rider/messaging"
	"github.com/EarvinKayonga/rider/storage"
)

// ListenerToBikeEvent for bike events.
func ListenerToBikeEvent(ctx context.Context, conf configuration.BikeConfiguration,
	logger logging.Logger, database storage.BikeStore) (messaging.Consumer, error) {

	listener, err := messaging.NewConsumer(ctx, conf.Messaging.Consumption, logger,
		func(message *bus.Message) (interface{}, error) {
			defer message.Finish()

			m := TrackTripPayload{}

			err := message.DecodePayload(&m)
			if err != nil {
				logger.
					WithError(err).
					Error("an error occured while decoding bike event payload")

				return nil, errors.Wrap(err,
					"an error occured while decoding bike event payload")
			}

			err = database.UpdateBikeLocation(ctx, m.BikeID, m.Lat, m.Lng)
			if err != nil {
				logger.
					WithError(err).
					Error("an error occured while updating bike location")

				return nil, errors.Wrap(err,
					"an error occured while updating bike location")
			}

			logger.Info("bike location successfully updated")

			return nil, nil
		})

	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while creating a consumer")
	}

	logger.Info("bike event consumer created")

	return listener, nil
}

// ListenerToTripEvent for trip events.
func ListenerToTripEvent(ctx context.Context, conf configuration.TripConfiguration,
	logger logging.Logger, database storage.TripStore) (messaging.Consumer, error) {

	listener, err := messaging.NewConsumer(ctx, conf.Messaging.Consumption, logger,
		func(message *bus.Message) (interface{}, error) {
			defer message.Finish()

			m := TrackTripPayload{}

			err := message.DecodePayload(&m)
			if err != nil {
				logger.
					WithError(err).
					Error("an error occured while decoding trip event payload")

				return nil, errors.Wrap(err,
					"an error occured while decoding trip event payload")
			}

			err = database.AddLocationToTrip(ctx, m.TripID, m.Lat, m.Lng)
			if err != nil {
				logger.
					WithError(err).
					Error("an error occured while updating trip with a location")

				return nil, errors.Wrap(err,
					"an error occured while updating a trip with a location")
			}

			logger.Info("trip location successfully updated")

			return nil, nil
		})

	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while creating a consumer")
	}

	logger.Info("trip event consumer created")

	return listener, nil
}
