package domain

import (
	"context"

	"github.com/EarvinKayonga/rider/models"
	"github.com/EarvinKayonga/rider/storage"
)

// StartTrip unsuprisingly starts a trip when possible.
func StartTrip(ctx context.Context, bikeID string,
	lat, lng float64) (*models.Trip, error) {
	return storage.TripStoreFromContext(ctx).CreateTrip(ctx, bikeID, lat, lng)
}

// EndTrip unsuprisingly starts a trip when possible.
func EndTrip(ctx context.Context, tripID string, lat, lng float64) (*models.Trip, error) {
	return storage.TripStoreFromContext(ctx).EndTrip(ctx, tripID, lat, lng)
}
