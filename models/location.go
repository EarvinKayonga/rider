package models

// Location model.
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
