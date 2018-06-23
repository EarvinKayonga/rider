package messaging

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	bus "github.com/rafaeljesus/nsq-event-bus"

	"github.com/EarvinKayonga/rider/configuration"
	"github.com/EarvinKayonga/rider/logging"
)

// Consumer  is an abstraction for consumming messages
// for messaging pipeline.
type Consumer interface {
	Run(ctx context.Context) error
}

type consumer struct {
	Address     string
	Topic       string
	HandlerFunc func(*bus.Message) (reply interface{}, err error)

	logger logging.Logger
}

func (e consumer) Run(ctx context.Context) error {
	errs := make(chan error, 1)

	go func() {
		err := bus.On(bus.ListenerConfig{
			Lookup:      []string{e.Address},
			Topic:       e.Topic,
			Channel:     fmt.Sprintf("consumer%d", rand.Intn(100)),
			HandlerFunc: e.HandlerFunc,
		})

		if err != nil {
			errs <- errors.Wrapf(err,
				"an error occured while setting consumer at %s", e.Address)
		}
	}()

	e.logger.Infof(
		"launching event consumer at %s, topic: %s, channel: %s",
		e.Address, e.Topic, "rider")

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		e.logger.Info("event consumer is shut down")
		return nil
	}
}

// NewConsumer returns a valid event consumer.
func NewConsumer(ctx context.Context, conf configuration.Consumption,
	logger logging.Logger,
	handler func(*bus.Message) (reply interface{}, err error)) (Consumer, error) {

	return &consumer{
		Address:     conf.Address,
		Topic:       conf.Topic,
		HandlerFunc: handler,

		logger: logger,
	}, nil
}
