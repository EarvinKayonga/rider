package entropy

import (
	"context"
	"io"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

type keyType string

const (
	key = keyType("entropy")
)

// IDGenerator specification.
type IDGenerator interface {
	NewID() string
	Compare(string, string) int
}

// FromContext extracts an IDGenerator from the Context.
func FromContext(ctx context.Context) IDGenerator {
	return ctx.Value(key).(IDGenerator)
}

// NewContext added an IDGenerator to the Context.
func NewContext(ctx context.Context, id IDGenerator) context.Context {
	return context.WithValue(ctx, key, id)
}

// NewIDGenerator creates an IDGenerator.
func NewIDGenerator() IDGenerator {
	t := time.Unix(1000000, 0)
	return &muon{
		entropy: rand.New(rand.NewSource(t.UnixNano())),
	}
}

type muon struct {
	entropy io.Reader
}

// NewID returns a new ID.
func (e *muon) NewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), e.entropy).String()
}

// NewID compares IDs.
func (e *muon) Compare(a, b string) int {
	return ulid.MustParse(a).Compare(ulid.MustParse(b))
}
