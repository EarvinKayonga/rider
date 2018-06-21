package storage

import (
	"context"

	"github.com/EarvinKayonga/rider/models"
)

// Store specifies  how persistance is handled in the app.
type Store interface {
	BikeStore
	TripStore

	Close(ctx context.Context) error
}

// BikeStore specifies how bike service persisted
// its data.
type BikeStore interface {
	CreateBikes(ctx context.Context, bikes []Bike) ([]models.Bike, error)
	ListBikes(ctx context.Context, cursor string, limit int64) ([]models.Bike, error)
	FindBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error)
	UpdateBikeLocation(ctx context.Context, bikeID string, lat, lng float64) error
	UnLockBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error)
	LockBikeByPublicID(ctx context.Context, bikeID string) (*models.Bike, error)
	ListAllBikes(ctx context.Context, limit int64) ([]models.Bike, error)
}

// TripStore specifies how trip service persisted
// its data.
type TripStore interface {
	GetLocationsForTrip(ctx context.Context, tripID string) ([]models.Location, error)
	AddLocationToTrip(ctx context.Context, tripID string, lat, lng float64) error
	CreateTrip(ctx context.Context, bikeID string, lat, lng float64) (*models.Trip, error)
	EndTrip(ctx context.Context, tripID string, lat, lng float64) (*models.Trip, error)
}
