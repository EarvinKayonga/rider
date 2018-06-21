package messaging

import (
	"context"
	"time"

	"github.com/pkg/errors"
	bus "github.com/rafaeljesus/nsq-event-bus"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/logging"
)

// Emitter sends message through messaging pipeline.
type Emitter interface {
	Emit(ctx context.Context, payload interface{}) error
}

type emitter struct {
	Emitter bus.Emitter
	Topic   string
}

// NewEmitter creates a valid emitter thrrough nsq.
func NewEmitter(ctx context.Context, conf configuration.Emission,
	logger logging.Logger) (Emitter, error) {
	emit, err := bus.NewEmitter(bus.EmitterConfig{
		DialTimeout:        1 * time.Second,
		ReadTimeout:        15 * time.Second,
		WriteTimeout:       4 * time.Second,
		MaxBackoffDuration: 30 * time.Second,
		MsgTimeout:         10 * time.Second,
		HeartbeatInterval:  5 * time.Second,

		MaxInFlight: conf.MaxInFlight,
		Address:     conf.Address,
	})

	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while contacting nsq with %s",
			conf.Address)
	}

	logger.Info("nsq emitter was created")

	return &emitter{
		Emitter: *emit,
		Topic:   conf.Topic,
	}, nil
}

func (e *emitter) Emit(_ context.Context, payload interface{}) error {
	return e.Emitter.EmitAsync(e.Topic, payload)
}
