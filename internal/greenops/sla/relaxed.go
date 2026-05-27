// Package sla expresses the "green SLA": low-priority paths are
// allowed relaxed latency targets so the scheduler can defer them to
// low-carbon hours.
package sla

import "time"

// Path is a named request class.
type Path string

const (
	PathSubmit      Path = "submit"
	PathView        Path = "view"
	PathAggregate   Path = "aggregate"
	PathBatchRedTeam Path = "batch-red-team"
)

// Target returns the latency target for the path.
func Target(p Path) time.Duration {
	switch p {
	case PathSubmit, PathView:
		return 8 * time.Second
	case PathAggregate:
		return 60 * time.Second
	case PathBatchRedTeam:
		return 6 * time.Hour
	}
	return 30 * time.Second
}

// EcoAllowed reports whether the green scheduler may defer this path.
func EcoAllowed(p Path) bool {
	switch p {
	case PathSubmit, PathView:
		return false
	}
	return true
}
