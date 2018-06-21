package domain

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/httpx"
	"github.com/EarvinKayonga/rider/models"
)

const (
	jsonContentType = "application/json"
)

// LockBikeFromGateway locks a bike.
func LockBikeFromGateway(ctx context.Context, conf configuration.GatewayConfiguration, bikeID string) (*models.Bike, error) {
	resp, err := httpx.Client().Get(conf.BikeURL + "/lock/" + bikeID)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while connecting bike service")
	}

	bike, err := DeserializeBikeFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return bike, nil
}

// UnLockBikeFromGateway unlocks a bike.
func UnLockBikeFromGateway(ctx context.Context, conf configuration.GatewayConfiguration, bikeID string) error {
	resp, err := httpx.Client().Get(conf.BikeURL + "/unlock/" + bikeID)
	if err != nil {
		return errors.Wrap(err, "an error occured while connecting bike service")
	}

	_, err = DeserializeBikeFromResponse(resp)
	return err
}

// DeserializeBikesFromResponse tries to extract an array of bikes for an http Response.
func DeserializeBikesFromResponse(resp *http.Response) ([]models.Bike, error) {
	if resp != nil {
		defer func() {
			thr := resp.Body.Close()
			_ = thr
		}()
	}

	if resp.Body == nil {
		return nil, ErrEmptyBody
	}

	bikes := []models.Bike{}
	err := json.NewDecoder(resp.Body).Decode(&bikes)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while decoding message from bike service")
	}

	return bikes, nil
}

// DeserializeBikeFromResponse tries to extract a bike for an http Response.
func DeserializeBikeFromResponse(resp *http.Response) (*models.Bike, error) {
	if resp != nil {
		defer func() {
			thr := resp.Body.Close()
			_ = thr
		}()
	}

	if resp.Body == nil {
		return nil, ErrEmptyBody
	}

	bike := models.Bike{}
	err := json.NewDecoder(resp.Body).Decode(&bike)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while decoding message from bike service")
	}

	return &bike, nil
}

// DeserializeTripFromResponse tries to extract a trip for an http Response.
func DeserializeTripFromResponse(resp *http.Response) (*models.Trip, error) {
	if resp != nil {
		defer func() {
			thr := resp.Body.Close()
			_ = thr
		}()
	}

	if resp.Body == nil {
		return nil, ErrEmptyBody
	}

	trip := models.Trip{}
	err := json.NewDecoder(resp.Body).Decode(&trip)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while decoding message from trip service")
	}

	return &trip, nil
}

// StartTripFromGateway calls the trip service to start a trip.
func StartTripFromGateway(ctx context.Context, conf configuration.GatewayConfiguration, bikeID string, lat, lng float64) (*models.Trip, error) {
	body, err := createStartTripBody(bikeID, lat, lng)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while create start trip payload")
	}

	resp, err := httpx.Client().Post(conf.BikeURL+"/trip/end", jsonContentType, body)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while connecting trip service")
	}
	if resp != nil {
		defer func() {
			thr := resp.Body.Close()
			_ = thr
		}()
	}

	if resp.Body == nil {
		return nil, ErrEmptyBody
	}

	trip := models.Trip{}
	err = json.NewDecoder(resp.Body).Decode(&trip)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while decoding message from bike service")
	}

	return &trip, nil
}

// EndTripFromGateway calls the trip service to end a trip.
func EndTripFromGateway(ctx context.Context, conf configuration.GatewayConfiguration, tripID string, lat, lng float64) (*models.Trip, error) {
	body, err := createEndTripBody(tripID, lat, lng)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while create end trip payload")
	}

	resp, err := httpx.Client().Post(conf.BikeURL+"/trip/end", jsonContentType, body)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while connecting trip service")
	}
	if resp != nil {
		defer func() {
			thr := resp.Body.Close()
			_ = thr
		}()
	}

	if resp.Body == nil {
		return nil, ErrEmptyBody
	}

	trip := models.Trip{}
	err = json.NewDecoder(resp.Body).Decode(&trip)
	if err != nil {
		return nil, errors.Wrap(err, "an error occured while decoding message from bike service")
	}

	return &trip, nil
}
