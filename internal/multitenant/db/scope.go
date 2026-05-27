// Package db wires the tenant scope into the data layer. Every query
// must run through Scope(ctx) so the tenant_id predicate is appended
// even when developers forget.
package db

import (
	"context"
	"errors"
)

type ctxKey struct{}

// WithTenant attaches the current tenant id to the context.
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, tenantID)
}

// FromContext returns the tenant id, or error.
func FromContext(ctx context.Context) (string, error) {
	if v, ok := ctx.Value(ctxKey{}).(string); ok && v != "" {
		return v, nil
	}
	return "", errors.New("no tenant on context")
}

// Optional crosses tenant boundaries for opt-in intelligence sharing.
// Callers must explicitly use this — default behaviour requires a
// tenant.
func Optional(ctx context.Context, share bool) (string, bool) {
	if share {
		return "", true
	}
	t, err := FromContext(ctx)
	if err != nil {
		return "", false
	}
	return t, false
}
