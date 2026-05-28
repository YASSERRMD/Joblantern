// Package effectiveness produces the monthly takedown-effectiveness
// report. It tracks how many packets were sent, accepted, and how
// long the targeted domains stayed reachable.
package effectiveness

import "time"

// Outcome enumerates resolutions per packet.
type Outcome string

const (
	OutcomeAccepted  Outcome = "accepted"
	OutcomeRejected  Outcome = "rejected"
	OutcomeIgnored   Outcome = "ignored"
	OutcomeRespawned Outcome = "respawned"
)

// Row is one packet's outcome row.
type Row struct {
	PacketID    string
	Domain      string
	Sent        time.Time
	Resolved    time.Time
	Outcome     Outcome
	HoursOnline int
}

// Summary is the monthly aggregate.
type Summary struct {
	Month             time.Month
	Year              int
	Total             int
	Accepted          int
	Rejected          int
	Ignored           int
	Respawned         int
	MedianHoursOnline int
}

// Aggregate produces the monthly summary.
func Aggregate(year int, m time.Month, rows []Row) Summary {
	s := Summary{Month: m, Year: year, Total: len(rows)}
	hours := []int{}
	for _, r := range rows {
		switch r.Outcome {
		case OutcomeAccepted:
			s.Accepted++
		case OutcomeRejected:
			s.Rejected++
		case OutcomeIgnored:
			s.Ignored++
		case OutcomeRespawned:
			s.Respawned++
		}
		hours = append(hours, r.HoursOnline)
	}
	if len(hours) > 0 {
		s.MedianHoursOnline = median(hours)
	}
	return s
}

func median(xs []int) int {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j-1] > xs[j]; j-- {
			xs[j-1], xs[j] = xs[j], xs[j-1]
		}
	}
	return xs[len(xs)/2]
}
