// Package conflict tracks council conflicts of interest.
package conflict

import "time"

// Declaration is a single COI disclosure.
type Declaration struct {
	SeatID       string
	DeclaredAt   time.Time
	Scope        string // free-text; e.g. "funded by ACME"
	Active       bool
}

// MustRecuse reports whether a seat should recuse from a case whose
// affected parties include affectedScopes.
func MustRecuse(decs []Declaration, affected []string) bool {
	scopes := map[string]bool{}
	for _, a := range affected {
		scopes[a] = true
	}
	for _, d := range decs {
		if d.Active && scopes[d.Scope] {
			return true
		}
	}
	return false
}
