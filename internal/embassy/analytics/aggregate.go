// Package analytics aggregates kiosk activity for ministry reporting.
// All metrics are anonymous and bucketed by week.
package analytics

import "time"

// WeekStats is one week of activity for one kiosk.
type WeekStats struct {
	KioskID            string
	WeekStart          time.Time
	Sessions           int
	GreenVerdicts      int
	YellowVerdicts     int
	RedVerdicts        int
	OverridesByOfficer int
	AvgSessionSeconds  int
}

// Roll computes a fleet-wide summary from per-kiosk stats. No
// per-applicant data crosses this boundary.
func Roll(in []WeekStats) WeekStats {
	out := WeekStats{}
	if len(in) == 0 {
		return out
	}
	out.WeekStart = in[0].WeekStart
	for _, s := range in {
		out.Sessions += s.Sessions
		out.GreenVerdicts += s.GreenVerdicts
		out.YellowVerdicts += s.YellowVerdicts
		out.RedVerdicts += s.RedVerdicts
		out.OverridesByOfficer += s.OverridesByOfficer
	}
	if out.Sessions > 0 {
		tot := 0
		for _, s := range in {
			tot += s.AvgSessionSeconds * s.Sessions
		}
		out.AvgSessionSeconds = tot / out.Sessions
	}
	return out
}
