// Package anonaudit reviews each archive for re-identification risk
// before publication.
package anonaudit

// Check is one audit row.
type Check struct {
	Name     string
	Severity string // "fatal" | "warning" | "info"
	Passed   bool
	Notes    string
}

// Default returns the canonical battery.
func Default() []Check {
	return []Check{
		{Name: "k-anonymity-k>=5", Severity: "fatal"},
		{Name: "no-raw-contacts", Severity: "fatal"},
		{Name: "no-raw-cvs", Severity: "fatal"},
		{Name: "stable-pseudonyms-rotated", Severity: "warning"},
		{Name: "geo-rounding-tile", Severity: "warning"},
	}
}

// Blocking returns true if any fatal check failed.
func Blocking(checks []Check) bool {
	for _, c := range checks {
		if c.Severity == "fatal" && !c.Passed {
			return true
		}
	}
	return false
}
