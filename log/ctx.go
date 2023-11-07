package log

import (
	"context"
	"io"
	"log/slog"
)

type _ctxKey struct{}

// ctxKey is the key used to store the logger in the context.
var ctxKey _ctxKey

// WithContext returns a new context with the given logger.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}

// FromContext returns the logger from the given context.
// If no logger is found, the default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// DefaultConsoleLogger returns a default console logger.
func DefaultConsoleLogger(w io.Writer, lvl slog.Leveler) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}))
}

// DefaultJSONLogger returns a default json logger.
func DefaultJSONLogger(w io.Writer, lvl slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}))
}
