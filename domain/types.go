package domain

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// EndTripPayload specifies the expected http body
// for ending a trip.
type EndTripPayload struct {
	TripID   string   `json:"trip_id"`
	Location Location `json:"location"`
}

// StartTripPayload specifies the expected http body
// for starting a trip.
type StartTripPayload struct {
	BikeID   string   `json:"bike_id"`
	Location Location `json:"location"`
}

// Location specifies location model.
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func createStartTripBody(bikeID string, lat, lng float64) (io.Reader, error) {
	start := StartTripPayload{
		BikeID: bikeID,
		Location: Location{
			Lat: lat,
			Lng: lng,
		},
	}

	jsonValue, err := json.Marshal(start)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while marshalling start trip payload")
	}

	return bytes.NewBuffer(jsonValue), nil
}

func createEndTripBody(tripID string, lat, lng float64) (io.Reader, error) {
	end := EndTripPayload{
		TripID: tripID,
		Location: Location{
			Lat: lat,
			Lng: lng,
		},
	}

	jsonValue, err := json.Marshal(end)
	if err != nil {
		return nil, errors.Wrap(err,
			"an error occured while marshalling end trip payload")
	}

	return bytes.NewBuffer(jsonValue), nil
}
