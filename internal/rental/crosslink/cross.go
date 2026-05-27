// Package crosslink detects when a job recruiter and a rental
// landlord share contact information — a strong signal that a single
// network is running both the employment and housing legs of the
// scam.
package crosslink

import "strings"

// Match is a detected overlap between two contexts.
type Match struct {
	Kind  string // "phone", "email", "domain", "address"
	Value string
}

// Job is the minimal recruiter contact view from the job verdict.
type Job struct {
	Phones    []string
	Emails    []string
	Domains   []string
	Addresses []string
}

// Rental is the minimal landlord contact view from the rental submission.
type Rental struct {
	Phones    []string
	Emails    []string
	Domains   []string
	Addresses []string
}

// Detect returns the overlaps between a job and a rental. Empty
// values on either side are ignored.
func Detect(j Job, r Rental) []Match {
	var out []Match
	out = append(out, intersect("phone", j.Phones, r.Phones)...)
	out = append(out, intersect("email", j.Emails, r.Emails)...)
	out = append(out, intersect("domain", j.Domains, r.Domains)...)
	out = append(out, intersect("address", j.Addresses, r.Addresses)...)
	return out
}

func intersect(kind string, a, b []string) []Match {
	set := map[string]bool{}
	for _, x := range a {
		if v := strings.ToLower(strings.TrimSpace(x)); v != "" {
			set[v] = true
		}
	}
	var out []Match
	seen := map[string]bool{}
	for _, x := range b {
		v := strings.ToLower(strings.TrimSpace(x))
		if v == "" || seen[v] || !set[v] {
			continue
		}
		seen[v] = true
		out = append(out, Match{Kind: kind, Value: v})
	}
	return out
}
