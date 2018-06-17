package models

// Bike model.
type Bike struct {
	ID       string   `json:"id"`
	Status   int      `json:"status"`
	Location Location `json:"location"`
}
