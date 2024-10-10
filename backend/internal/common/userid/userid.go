// Package userid provides a context-based user ID.
package userid

// Inspired by https://go.dev/blog/context

import (
	"context"
)

// contextKey is a context key for the GitHub user ID. We retrieve it from the JWT in incoming requests.
type contextKey struct{}

func WithUserID(ctx context.Context, jwt int64) context.Context {
	return context.WithValue(ctx, contextKey{}, jwt)
}

func FromContext(ctx context.Context) (int64, bool) {
	jwt, ok := ctx.Value(contextKey{}).(int64)
	return jwt, ok
}
