// Package routing maps a verdict context (origin country, destination
// country, industry) to the regulators that should receive
// notifications. The matrix is data-driven so new jurisdictions can
// be added without code changes.
package routing

import "strings"

// Match represents one regulator subscription to a context.
type Match struct {
	Regulator      string
	Country        string
	IndustryFilter []string
}

// Table is the in-memory routing table.
type Table struct {
	entries []Match
}

// NewTable returns an empty routing table.
func NewTable() *Table { return &Table{} }

// Add adds a subscription.
func (t *Table) Add(m Match) { t.entries = append(t.entries, m) }

// Resolve returns the regulators that match the supplied verdict
// context. A match requires the country to align AND, if an industry
// filter is set, the verdict industry to be in it.
func (t *Table) Resolve(country, industry string) []string {
	country = strings.ToUpper(country)
	industry = strings.ToLower(industry)
	var out []string
	seen := map[string]bool{}
	for _, e := range t.entries {
		if strings.ToUpper(e.Country) != country {
			continue
		}
		if len(e.IndustryFilter) > 0 {
			ok := false
			for _, f := range e.IndustryFilter {
				if strings.ToLower(f) == industry {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if !seen[e.Regulator] {
			seen[e.Regulator] = true
			out = append(out, e.Regulator)
		}
	}
	return out
}
