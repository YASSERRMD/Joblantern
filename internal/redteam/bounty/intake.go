// Package bounty is the scaffolding for external bug-bounty intake.
// We accept reports through a stable endpoint and route them into the
// Trust & Safety council case docket.
package bounty

import "time"

// Submission is the wire form of a bounty submission.
type Submission struct {
	ID             string
	ReporterEmail  string
	Severity       string // "low", "medium", "high", "critical"
	Title          string
	Summary        string
	ReceivedAt     time.Time
	PgpFingerprint string
}

// Severity ranks the report for routing.
var SeverityOrder = map[string]int{"low": 0, "medium": 1, "high": 2, "critical": 3}

// Triage assigns an initial response window.
func Triage(s Submission) time.Duration {
	switch s.Severity {
	case "critical":
		return 2 * time.Hour
	case "high":
		return 24 * time.Hour
	case "medium":
		return 72 * time.Hour
	}
	return 7 * 24 * time.Hour
}
