package logging

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/EarvinKayonga/rider/configuration"
)

// keyType is a package local type alias purposed to avoid name collision
// in context.
type keyType string

const (
	key = keyType("logging")
)

// NewLogger creates a logger instance.
func NewLogger(conf configuration.Logging) *Logger {
	lvl, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		lvl = logrus.WarnLevel
	}

	logger := logrus.New()
	logger.Level = lvl
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	if conf.Format == configuration.JSONFormat {
		logger.Formatter = &logrus.JSONFormatter{}
	}

	return &Logger{
		logger,
	}
}

// Logger is a logging abstraction.
type Logger struct {
	*logrus.Logger
}

// FromContext extracts a Logger from the Context.
func FromContext(ctx context.Context) *Logger {
	return ctx.Value(key).(*Logger)
}

// NewContext adds the given Logger to the Context.
func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, key, l)
}
