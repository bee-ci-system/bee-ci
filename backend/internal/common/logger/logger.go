package logger

// Inspired by https://go.dev/blog/context

import (
	"context"
	"log/slog"
)

type contextKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func FromContext(ctx context.Context) (*slog.Logger, bool) {
	logger, ok := ctx.Value(contextKey{}).(*slog.Logger)
	return logger, ok
}
