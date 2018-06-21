package stats

import (
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/pkg/errors"

	"github.com/EarvinKayonga/rider/configuration"
)

// Statter is an abstraction for handling
// stats on the running app.
// A Mock implementation is provided.
type Statter interface {
	Close() error
	Dec(stat string, value int64, rate float32) error
	Gauge(stat string, value int64, rate float32) error
	GaugeDelta(stat string, value int64, rate float32) error
	Inc(stat string, value int64, rate float32) error
	Raw(stat string, value string, rate float32) error
	Set(stat string, value string, rate float32) error
	SetInt(stat string, value int64, rate float32) error
	SetPrefix(prefix string)
	Timing(stat string, delta int64, rate float32) error
	TimingDuration(stat string, delta time.Duration, rate float32) error
}

// statter is a concret implementation for statter
// using github.com/cactus/go-statsd-client.BufferedClient.
type statter struct {
	statsd.Statter
}

// NewStatsdClient returns a valid statsd client.
func NewStatsdClient(conf configuration.Monitoring) (Statter, error) {
	// flushInterval is a time.Duration, and specifies the maximum interval for
	// packet sending. Note that if you send lots of metrics, you will send more
	// often. This is just a maximal threshold.
	flushInterval := 300 * time.Millisecond

	// If flushBytes is 0, defaults to 1432 bytes, which is considered safe
	// for local traffic. If sending over the public internet, 512 bytes is
	// the recommended value.
	flushBytes := 512

	client, err := statsd.NewBufferedClient(conf.Addr, conf.Prefix, flushInterval, flushBytes)
	if err != nil {
		return nil, errors.Wrapf(err,
			"an error occured while creating statsd client at %s", conf.Addr)
	}

	return &statter{
		client,
	}, nil
}
