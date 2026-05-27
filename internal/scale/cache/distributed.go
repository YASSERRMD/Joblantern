// Package cache documents and abstracts the distributed cache layer.
// We support both Redis (BSD-3) and DragonflyDB (Apache 2.0). The
// interface here is the small subset Joblantern actually uses.
package cache

import (
	"context"
	"errors"
	"time"
)

// Cache is the abstraction.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

// ErrMiss is returned on cache misses.
var ErrMiss = errors.New("cache miss")
