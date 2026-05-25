// Package cache is a tiny in-memory TTL cache used by MCP servers to
// avoid hammering rate-limited upstream APIs (Nominatim, Overpass,
// Mapillary, OpenCorporates, etc).
//
// The implementation is intentionally minimal: a sync.Map of entries
// with a wall-clock expiry. For v1 scale this is fine; if a real
// shared cache is needed later we will swap in Redis behind the same
// interface.
package cache

import (
	"sync"
	"time"
)

// TTL is a single-process map cache with per-entry expiry.
type TTL[K comparable, V any] struct {
	mu     sync.RWMutex
	items  map[K]entry[V]
	maxAge time.Duration
}

type entry[V any] struct {
	value V
	exp   time.Time
}

// New returns a TTL cache that drops entries older than maxAge.
func New[K comparable, V any](maxAge time.Duration) *TTL[K, V] {
	return &TTL[K, V]{
		items:  make(map[K]entry[V]),
		maxAge: maxAge,
	}
}

// Get returns the cached value (and true) if present and unexpired.
func (c *TTL[K, V]) Get(key K) (V, bool) {
	var zero V
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return zero, false
	}
	if time.Now().After(e.exp) {
		c.mu.Lock()
		// Re-check under write lock in case another goroutine refreshed.
		if e2, still := c.items[key]; still && time.Now().After(e2.exp) {
			delete(c.items, key)
		}
		c.mu.Unlock()
		return zero, false
	}
	return e.value, true
}

// Set inserts or replaces the value for the given key.
func (c *TTL[K, V]) Set(key K, value V) {
	c.mu.Lock()
	c.items[key] = entry[V]{value: value, exp: time.Now().Add(c.maxAge)}
	c.mu.Unlock()
}

// Len returns the cache size, including expired entries that have not
// yet been swept.
func (c *TTL[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
