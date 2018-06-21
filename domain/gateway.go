package domain

import (
	"context"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/httpx"
	"github.com/EarvinKayonga/rider/models"
)

// GatewayListOfBikes lists bikes with pagination.
func GatewayListOfBikes(ctx context.Context, conf configuration.GatewayConfiguration, cursor string, limit int64) ([]models.Bike, error) {
	resp, err := httpx.Client().Get(conf.BikeURL + "/bikes/")
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while connecting bike service")
	}

	bikes, err := DeserializeBikesFromResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while deserializing bikes from http response")
	}

	return bikes, nil
}

// GatewayGetBikeByID returns a bike given an valid ID.
func GatewayGetBikeByID(ctx context.Context, conf configuration.GatewayConfiguration, bikeID string) (*models.Bike, error) {
	resp, err := httpx.Client().Get(conf.BikeURL + "/bike/" + bikeID)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while connecting bike service")
	}

	bike, err := DeserializeBikeFromResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while deserializing a bike from http reponse")
	}

	return bike, nil
}

// GatewayStartTrip unsuprisingly starts a trip when possible.
func GatewayStartTrip(ctx context.Context, conf configuration.GatewayConfiguration, bikeID string,
	lat, lng float64) (*models.Trip, error) {

	bike, err := LockBikeFromGateway(ctx, conf, bikeID)
	if err != nil {
		return nil, err
	}

	trip, err := StartTripFromGateway(ctx, conf, bike.ID, lat, lng)
	if err != nil {
		return nil, err
	}

	return trip, nil
}

// GatewayEndTrip unsuprisingly starts a trip when possible.
func GatewayEndTrip(ctx context.Context, conf configuration.GatewayConfiguration, tripID string, lat, lng float64) (*models.Trip, error) {
	trip, err := EndTripFromGateway(ctx, conf, tripID, lat, lng)
	if err != nil {
		return nil, err
	}

	err = UnLockBikeFromGateway(ctx, conf, trip.BikeID)
	if err != nil {
		return nil, err
	}

	return trip, nil
}
