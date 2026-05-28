// Package membership manages council seat lifecycle.
package membership

import (
	"errors"
	"time"
)

// SeatKind enumerates the seat types from the charter.
type SeatKind string

const (
	SeatNGO         SeatKind = "ngo"
	SeatTechnical   SeatKind = "technical"
	SeatAcademic    SeatKind = "academic"
	SeatRegulator   SeatKind = "regulator"
	SeatEngineering SeatKind = "engineering"
)

// Seat is one occupied seat.
type Seat struct {
	ID          string
	Kind        SeatKind
	Holder      string
	Affiliation string
	StartedAt   time.Time
	EndsAt      time.Time
	Suspended   bool
}

// Active reports whether the seat is currently held.
func (s Seat) Active(now time.Time) bool {
	if s.Suspended {
		return false
	}
	if !s.StartedAt.IsZero() && now.Before(s.StartedAt) {
		return false
	}
	if !s.EndsAt.IsZero() && now.After(s.EndsAt) {
		return false
	}
	return true
}

// EnsureQuorum returns nil if the supplied seats satisfy the
// charter's minimum composition.
func EnsureQuorum(seats []Seat) error {
	counts := map[SeatKind]int{}
	for _, s := range seats {
		if s.Active(time.Now()) {
			counts[s.Kind]++
		}
	}
	if counts[SeatNGO] < 3 {
		return errors.New("not enough active NGO seats")
	}
	if counts[SeatTechnical] < 2 {
		return errors.New("not enough technical seats")
	}
	if counts[SeatAcademic] < 2 {
		return errors.New("not enough academic seats")
	}
	return nil
}
