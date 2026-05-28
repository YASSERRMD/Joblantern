// Package privacy enforces the in-memory-only contract for
// personalization. CV bytes and parsed profiles live in process
// memory for the duration of a request and are zeroed on return.
//
// Persistence requires explicit user opt-in stored separately on the
// submission record.
package privacy

import "sync"

// Pool reuses byte slices so we can scrub them deterministically when
// they leave scope.
var Pool = sync.Pool{New: func() any { return make([]byte, 0, 4096) }}

// Scrub overwrites the slice with zeros and returns it to the pool.
func Scrub(b []byte) {
	for i := range b {
		b[i] = 0
	}
	b = b[:0]
	Pool.Put(b)
}

// Consent captures the user choice. Default is RetainNo.
type Consent string

const (
	RetainNo           Consent = "no"
	RetainResearchOnly Consent = "research-only"
	RetainImproveAgent Consent = "improve-agent"
)

// PersistAllowed returns whether the supplied consent permits writing
// the parsed profile to disk.
func PersistAllowed(c Consent) bool {
	return c == RetainResearchOnly || c == RetainImproveAgent
}
