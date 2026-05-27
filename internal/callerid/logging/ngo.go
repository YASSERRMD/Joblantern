// Package logging holds the opt-in detailed log NGOs use to
// investigate scam networks. Logs include call timing and outcomes,
// never recording or audio.
package logging

import "time"

// Event is one log entry.
type Event struct {
	At        time.Time `json:"at"`
	PhoneHash string    `json:"phone_hash"`
	Band      string    `json:"band"`
	Outcome   string    `json:"outcome"` // "rejected", "answered", "missed"
	DurationS int       `json:"duration_s,omitempty"`
}

// Consent is the user-side toggle gating this log.
type Consent struct {
	Enabled    bool
	NgoPartner string
	GrantedAt  time.Time
	Until      time.Time
}

// Active returns true if the consent is currently in force.
func (c Consent) Active(now time.Time) bool {
	if !c.Enabled {
		return false
	}
	if c.Until.IsZero() {
		return true
	}
	return now.Before(c.Until)
}
