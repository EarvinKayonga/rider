package domain

import (
	"context"

	"github.com/EarvinKayonga/rider/models"
	"github.com/EarvinKayonga/rider/storage"
)

// ListOfBikes lists bikes with pagination.
func ListOfBikes(ctx context.Context, cursor string, limit int64) ([]models.Bike, error) {
	return storage.
		BikeStoreFromContext(ctx).
		ListBikes(ctx, cursor, limit)
}

// GetBikeByID returns a bike given an valid ID.
func GetBikeByID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return storage.
		BikeStoreFromContext(ctx).
		FindBikeByPublicID(ctx, bikeID)
}

// LockBikeByID locks a bike given an valid ID.
func LockBikeByID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return storage.
		BikeStoreFromContext(ctx).
		LockBikeByPublicID(ctx, bikeID)
}

// UnLockBikeByID unlocks a bike given an valid ID.
func UnLockBikeByID(ctx context.Context, bikeID string) (*models.Bike, error) {
	return storage.
		BikeStoreFromContext(ctx).
		UnLockBikeByPublicID(ctx, bikeID)
}
