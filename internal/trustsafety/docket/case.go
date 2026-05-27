// Package docket is the council's case docket. Every appeal,
// emergency takedown, or charter amendment becomes a docketed case.
package docket

import "time"

// Kind enumerates docket case classes.
type Kind string

const (
	KindAppeal           Kind = "appeal"
	KindEmergency        Kind = "emergency-takedown"
	KindCharterAmend     Kind = "charter-amendment"
	KindMembershipChange Kind = "membership-change"
)

// Case is one docketed case.
type Case struct {
	ID         string
	Kind       Kind
	Title      string
	OpenedAt   time.Time
	ClosedAt   time.Time
	Owner      string
	Status     string
	PublicMins bool
}

// Open reports whether the case is open.
func (c Case) Open() bool { return c.ClosedAt.IsZero() }
