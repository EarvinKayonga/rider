package domain

import (
	"context"

	"github.com/EarvinKayonga/rider/messaging"
)

// TrackTripPayload for messaging.
type TrackTripPayload struct {
	Lng, Lat float64
	BikeID   string
	TripID   string
}

// TrackTrip sends a payload through the messaging pipeline.
func TrackTrip(ctx context.Context,
	messenger messaging.Emitter, hearbeat TrackTripPayload) error {

	return messenger.Emit(ctx, &hearbeat)
}
