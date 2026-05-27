// Package disclosure documents the responsible-disclosure pipeline.
package disclosure

import "time"

// Policy captures the operational windows.
type Policy struct {
	AcknowledgmentWindow time.Duration
	FixWindowCritical    time.Duration
	FixWindowHigh        time.Duration
	PublicDisclosure     time.Duration
}

// Default is the published policy.
func Default() Policy {
	return Policy{
		AcknowledgmentWindow: 48 * time.Hour,
		FixWindowCritical:    7 * 24 * time.Hour,
		FixWindowHigh:        30 * 24 * time.Hour,
		PublicDisclosure:     90 * 24 * time.Hour,
	}
}
