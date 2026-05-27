// Package fleet manages a deployed kiosk fleet — heartbeat, version
// inventory, alerting on stale kiosks.
package fleet

import (
	"sync"
	"time"
)

// Kiosk is one deployed device.
type Kiosk struct {
	ID            string
	Embassy       string
	Country       string
	Version       string
	LastHeartbeat time.Time
	OnlineHours24h int
}

// Fleet is the in-memory inventory.
type Fleet struct {
	mu      sync.Mutex
	kiosks  map[string]Kiosk
	staleAt time.Duration
}

// New returns a fleet whose staleness threshold is t.
func New(staleAfter time.Duration) *Fleet {
	return &Fleet{kiosks: map[string]Kiosk{}, staleAt: staleAfter}
}

// Heartbeat records a kiosk check-in.
func (f *Fleet) Heartbeat(k Kiosk) {
	f.mu.Lock()
	defer f.mu.Unlock()
	k.LastHeartbeat = time.Now().UTC()
	f.kiosks[k.ID] = k
}

// Stale returns kiosks that have not checked in within the threshold.
func (f *Fleet) Stale() []Kiosk {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []Kiosk
	cutoff := time.Now().UTC().Add(-f.staleAt)
	for _, k := range f.kiosks {
		if k.LastHeartbeat.Before(cutoff) {
			out = append(out, k)
		}
	}
	return out
}

// All returns a snapshot of all kiosks.
func (f *Fleet) All() []Kiosk {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]Kiosk, 0, len(f.kiosks))
	for _, k := range f.kiosks {
		out = append(out, k)
	}
	return out
}
