// Package crosslink ties a marketplace submission back to the jobs
// and rental graphs. A network that runs all three is a clear
// scam-operator signature.
package crosslink

import "strings"

// Match is one detected overlap.
type Match struct {
	Kind  string // "phone", "email"
	Value string
}

// Detect returns overlaps across the supplied identifier sets.
func Detect(jobPhones, rentalPhones, mktPhones []string) []Match {
	set := func(xs []string) map[string]bool {
		m := map[string]bool{}
		for _, x := range xs {
			if v := strings.TrimSpace(x); v != "" {
				m[v] = true
			}
		}
		return m
	}
	jp := set(jobPhones)
	rp := set(rentalPhones)
	var out []Match
	seen := map[string]bool{}
	for _, p := range mktPhones {
		v := strings.TrimSpace(p)
		if v == "" || seen[v] {
			continue
		}
		if jp[v] || rp[v] {
			seen[v] = true
			out = append(out, Match{Kind: "phone", Value: v})
		}
	}
	return out
}
