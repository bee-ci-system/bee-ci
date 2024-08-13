package logs

import (
	"context"
	"log/slog"
)

type ContextKey struct{}

func FromContext(ctx context.Context) (*slog.Logger, bool) {
	logger, ok := ctx.Value(ContextKey{}).(*slog.Logger)
	return logger, ok
}
