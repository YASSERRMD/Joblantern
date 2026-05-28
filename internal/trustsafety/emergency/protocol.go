// Package emergency implements the fast-track takedown protocol.
//
// A 3-of-5 emergency vote can authorise a takedown packet within the
// hour; the full council must ratify within 72 hours or the action
// is rolled back and the decision is logged.
package emergency

import (
	"errors"
	"time"
)

// Vote is one emergency vote.
type Vote struct {
	SeatID   string
	CastAt   time.Time
	Decision string // "approve" / "deny"
}

// Action is the emergency case in progress.
type Action struct {
	ID         string
	OpenedAt   time.Time
	Votes      []Vote
	Ratified   bool
	Rolledback bool
}

// Threshold is 3 approvals.
const Threshold = 3

// Approved returns true if the action has reached the threshold and
// has not been rolled back.
func (a Action) Approved() bool {
	if a.Rolledback {
		return false
	}
	yes := 0
	for _, v := range a.Votes {
		if v.Decision == "approve" {
			yes++
		}
	}
	return yes >= Threshold
}

// MustRatifyBy returns the deadline for full-council ratification.
func (a Action) MustRatifyBy() time.Time { return a.OpenedAt.Add(72 * time.Hour) }

// FailIfMissed returns an error if the deadline has passed without
// ratification.
func FailIfMissed(a Action, now time.Time) error {
	if a.Ratified {
		return nil
	}
	if now.After(a.MustRatifyBy()) {
		return errors.New("emergency action lapsed without ratification")
	}
	return nil
}
