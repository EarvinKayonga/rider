package models

// Location model.
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// CreateLocation creates a location from given latitude and longitude.
func CreateLocation(lat, lng float64) Location {
	return Location{
		Type: "Point",
		Coordinates: []float64{
			lat,
			lng,
		},
	}
}
