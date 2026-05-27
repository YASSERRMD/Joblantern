// Package cache stores a per-region snapshot of high-confidence
// scam-number hashes so the device can serve a verdict in under a
// second even when the network is unavailable.
package cache

import "sync"

// Snapshot is a region-scoped cache.
type Snapshot struct {
	mu     sync.RWMutex
	red    map[string]struct{}
	yellow map[string]struct{}
	green  map[string]struct{}
}

// New returns an empty snapshot.
func New() *Snapshot {
	return &Snapshot{
		red:    map[string]struct{}{},
		yellow: map[string]struct{}{},
		green:  map[string]struct{}{},
	}
}

// Replace swaps in new sets atomically.
func (s *Snapshot) Replace(red, yellow, green []string) {
	r := setOf(red)
	y := setOf(yellow)
	g := setOf(green)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.red, s.yellow, s.green = r, y, g
}

// Band returns the band for the supplied hash, or empty string if unknown.
func (s *Snapshot) Band(hash string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.red[hash]; ok {
		return "red"
	}
	if _, ok := s.yellow[hash]; ok {
		return "yellow"
	}
	if _, ok := s.green[hash]; ok {
		return "green"
	}
	return ""
}

func setOf(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}
