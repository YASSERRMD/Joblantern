// Package offline gives the kiosk a local-first cache so it can keep
// scoring submissions during connectivity gaps. The cache is a tiny
// embedded KV that is replenished when connectivity returns and that
// flushes any locally produced verdicts back to the central server.
package offline

import (
	"sync"
	"time"
)

// Entry is one cached lookup keyed by document hash.
type Entry struct {
	Key       string
	Value     []byte
	StoredAt  time.Time
	ExpiresAt time.Time
}

// Cache is an in-memory LRU-style cache. Production implementations
// swap in BoltDB or SQLite for durability across kiosk reboots.
type Cache struct {
	mu   sync.Mutex
	data map[string]Entry
	cap  int
}

// New creates a cache with capacity cap.
func New(cap int) *Cache { return &Cache{data: map[string]Entry{}, cap: cap} }

// Put stores an entry.
func (c *Cache) Put(e Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.data) >= c.cap {
		// drop one arbitrary entry — production uses LRU
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}
	c.data[e.Key] = e
}

// Get returns a fresh entry or false.
func (c *Cache) Get(key string) (Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[key]
	if !ok {
		return Entry{}, false
	}
	if !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt) {
		delete(c.data, key)
		return Entry{}, false
	}
	return e, true
}

// Pending returns the entries that have not been synced upstream
// since their StoredAt timestamp.
func (c *Cache) Pending(syncedThrough time.Time) []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []Entry
	for _, e := range c.data {
		if e.StoredAt.After(syncedThrough) {
			out = append(out, e)
		}
	}
	return out
}
