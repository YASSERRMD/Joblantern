// Package whitelist accepts official rosters of licensed recruiters
// from a regulator. A whitelist hit doesn't suppress the verdict but
// reduces baseline risk and pulls in the regulator's licence-status
// metadata.
package whitelist

import "time"

// Entry is one licensed recruiter on a regulator whitelist.
type Entry struct {
	RegulatorID  string
	EntityName   string
	EntityID     string
	LicenceID    string
	LicenceStart time.Time
	LicenceEnd   time.Time
	Scope        string
	URL          string
}

// RiskAdjustment is the absolute amount subtracted from the baseline
// risk score (out of 100) when a whitelist entry matches. The
// adjustment never crosses zero.
const RiskAdjustment = 20

// IsActive reports whether the licence covers the supplied moment.
func (e Entry) IsActive(at time.Time) bool {
	if e.LicenceStart.IsZero() && e.LicenceEnd.IsZero() {
		return false
	}
	if !e.LicenceStart.IsZero() && at.Before(e.LicenceStart) {
		return false
	}
	if !e.LicenceEnd.IsZero() && at.After(e.LicenceEnd) {
		return false
	}
	return true
}
