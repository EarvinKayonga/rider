package storage

import "github.com/pkg/errors"

// Database based errors.
var (
	ErrBikeNotFound   = errors.New("bike not found")
	ErrNotImplemented = errors.New("not implemented")
)
