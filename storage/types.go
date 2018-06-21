package storage

import (
	"time"

	"github.com/lib/pq"
)

// Location is the database representation of a models.Location.
type Location struct {
	ID        int64
	Latitude  float64
	Longitude float64
	TripID    string
	CreatedAt time.Time
}

// Trip is the database representation of a models.Trip.
type Trip struct {
	ID        int64
	PublicID  string
	Status    int
	BikeID    string
	StartedAt time.Time
	EndedAt   pq.NullTime
}

// Bike is the database representation of a models.Bike.
type Bike struct {
	ID        int64
	PublicID  string
	Status    int
	Latitude  float64
	Longitude float64
}
