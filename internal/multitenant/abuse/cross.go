// Package abuse detects cross-tenant abuse — e.g. a tenant using the
// service to enumerate or de-anonymise another tenant's verdicts.
// Detection runs as a passive analysis pass; remediation is council-
// driven.
package abuse

// Event is one suspicious cross-tenant lookup.
type Event struct {
	ActorTenant string
	TargetField string // "phone", "email", ...
	Window      string // "1h", "24h", ...
	Hits        int
}

// Threshold is the per-window hit count that flips a lookup pattern
// from "exploratory" to "suspect enumeration".
var Threshold = map[string]int{
	"1h":  500,
	"24h": 10000,
}

// IsAbusive returns true if the event crosses the configured threshold.
func IsAbusive(e Event) bool {
	if t, ok := Threshold[e.Window]; ok {
		return e.Hits > t
	}
	return false
}
