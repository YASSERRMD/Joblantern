// Package dpdpa implements the India DPDPA consent flow. Key
// requirements: explicit notice in English + the user's chosen
// regional language, granular per-purpose consent, easy withdrawal.
package dpdpa

import "time"

// Purpose enumerates each lawful purpose for processing.
type Purpose string

const (
	PurposeVerdict      Purpose = "verdict-generation"
	PurposeResearch     Purpose = "research-sharing"
	PurposeRegulator    Purpose = "regulator-forwarding"
	PurposeImproveAgent Purpose = "improve-agent"
)

// Consent is a per-purpose record.
type Consent struct {
	UserID    string
	Purpose   Purpose
	GrantedAt time.Time
	Notice    string // sha256 of the notice version shown
	Withdrawn time.Time
}

// Active reports whether the consent is currently in force.
func (c Consent) Active(now time.Time) bool {
	if c.GrantedAt.IsZero() {
		return false
	}
	if !c.Withdrawn.IsZero() && !now.Before(c.Withdrawn) {
		return false
	}
	return true
}
