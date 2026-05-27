// Package audit augments the existing audit log with ISO 27001
// retention and integrity-check requirements.
package audit

import "time"

// Retention is the per-class retention window.
type Retention struct {
	Class    string
	Duration time.Duration
}

// Defaults captures the policy.
func Defaults() []Retention {
	return []Retention{
		{Class: "access", Duration: 365 * 24 * time.Hour},
		{Class: "admin", Duration: 7 * 365 * 24 * time.Hour},
		{Class: "data-export", Duration: 7 * 365 * 24 * time.Hour},
		{Class: "incident", Duration: 7 * 365 * 24 * time.Hour},
	}
}

// IntegrityCheck is the daily Merkle-root verification expected by
// ISO 27001 controls. The verifier walks the chain and recomputes
// hashes.
type IntegrityCheck struct {
	At              time.Time
	EntriesChecked  int
	BrokenAt        int // -1 if intact
}
