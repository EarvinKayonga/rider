package models

import (
	"time"
)

// Trip model.
type Trip struct {
	ID        string     `json:"id"`
	Status    int        `json:"status"`
	BikeID    string     `json:"bike_id"`
	Locations []Location `json:"locations"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}
