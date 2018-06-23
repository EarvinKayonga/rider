package domain

import (
	"github.com/pkg/errors"
)

// Domain based errors.
var (
	ErrBikeInUse  = errors.New("bike already in use")
	ErrEmptyBody  = errors.New("empty body")
	ErrUnexpected = errors.New("unexpected error")
)
