// Package forwarding gates the optional pipeline that forwards
// confirmed-scam verdicts to the user's destination-country regulator.
// Forwarding is OFF by default and requires explicit per-submission
// user consent.
package forwarding

import (
	"errors"
	"time"
)

// Threshold is the minimum confidence at which forwarding is even
// offered to the user.
const Threshold = 0.85

// Request is what gets enqueued for forwarding.
type Request struct {
	VerdictID       string
	UserConsentedAt time.Time
	Regulator       string
	Country         string
	Confidence      float64
	Summary         string
}

// Validate returns nil if the request is eligible for forwarding.
func (r Request) Validate() error {
	if r.UserConsentedAt.IsZero() {
		return errors.New("user consent missing")
	}
	if r.Confidence < Threshold {
		return errors.New("verdict confidence below threshold")
	}
	if r.Regulator == "" || r.Country == "" {
		return errors.New("regulator routing incomplete")
	}
	return nil
}
