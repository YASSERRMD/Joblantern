// Package metrics defines the per-partner success metrics. The
// intent is not to grade partners but to surface where central
// support time should go.
package metrics

import "time"

// Snapshot is one per-partner measurement.
type Snapshot struct {
	PartnerID          string
	At                 time.Time
	VerdictsLast30d    int
	UniqueUsersLast30d int
	UnresolvedAppeals  int
	MentorshipMinutes  int
}

// Healthy reports whether the partner is hitting the soft floors. A
// partner can be unhealthy and still highly valuable — the function
// is a routing hint, not a verdict.
func Healthy(s Snapshot) bool {
	return s.VerdictsLast30d > 10 && s.UnresolvedAppeals < 5
}
