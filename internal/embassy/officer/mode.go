// Package officer is the consular-officer surface of the kiosk.
// Officers authenticate with a smart-card / staff badge and can:
//   - override a verdict with a documented reason,
//   - attach private counselling notes,
//   - reprint a previous receipt.
package officer

import (
	"errors"
	"strings"
	"time"
)

// Note is a private counselling note attached to a kiosk session.
type Note struct {
	OfficerID string
	At        time.Time
	Body      string
	Sealed    bool
}

// Override records an officer-issued change to a verdict band.
type Override struct {
	OfficerID  string
	At         time.Time
	FromBand   string
	ToBand     string
	Reason     string
	VerdictID  string
}

// ValidateOverride enforces that an officer override has a documented
// reason at least 30 characters long. Empty overrides are blocked.
func ValidateOverride(o Override) error {
	if strings.TrimSpace(o.OfficerID) == "" {
		return errors.New("officer id required")
	}
	if strings.TrimSpace(o.VerdictID) == "" {
		return errors.New("verdict id required")
	}
	if len([]rune(strings.TrimSpace(o.Reason))) < 30 {
		return errors.New("reason must be at least 30 characters")
	}
	if o.FromBand == o.ToBand {
		return errors.New("no-op override")
	}
	return nil
}
