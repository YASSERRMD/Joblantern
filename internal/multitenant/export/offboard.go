// Package export handles tenant export and offboarding. A tenant can
// pull a full machine-readable archive of their data at any time.
// Offboarding triggers deletion within 30 days, surfaced through the
// Phase 48 compliance pack.
package export

import "time"

// Archive describes the offboarding bundle.
type Archive struct {
	TenantID      string
	GeneratedAt   time.Time
	Verdicts      int
	Evidence      int
	BlocklistRow  int
	BytesEstimate int64
	URL           string
}

// OffboardSchedule returns the deletion deadline (now + 30 days).
func OffboardSchedule(now time.Time) time.Time {
	return now.Add(30 * 24 * time.Hour)
}
