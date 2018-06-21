package storage

import "context"

const (
	key = keyType("database")
)

type keyType string

// BikeStoreFromContext extracts a BikeStore from the Context.
func BikeStoreFromContext(ctx context.Context) BikeStore {
	return ctx.Value(key).(BikeStore)
}

// TripStoreFromContext extracts a TripStore from the Context.
func TripStoreFromContext(ctx context.Context) TripStore {
	return ctx.Value(key).(TripStore)
}

// NewContext adds the given Store to the Context.
func NewContext(ctx context.Context, db Store) context.Context {
	return context.WithValue(ctx, key, db)
}
