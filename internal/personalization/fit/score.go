// Package fit scores how well a candidate profile aligns with the
// role a listing claims to recruit for.
package fit

import "strings"

// Listing is a minimal view of the job-listing the user submitted.
type Listing struct {
	Role        string
	Industry    string
	Description string
}

// Profile is the candidate view derived from the optional CV.
type Profile struct {
	Role            string
	YearsExperience int
	Skills          []string
}

// Score returns a 0..100 value: how well the candidate's role fits
// the listing's role and surrounding industry context. Higher is a
// better match.
func Score(l Listing, p Profile) int {
	if l.Role == "" || p.Role == "" {
		return 50 // neutral when we cannot compare
	}
	score := 0
	if strings.EqualFold(strings.TrimSpace(l.Role), strings.TrimSpace(p.Role)) {
		score += 60
	} else if related(l.Role, p.Role) {
		score += 30
	}
	desc := strings.ToLower(l.Description)
	for _, s := range p.Skills {
		if strings.Contains(desc, strings.ToLower(s)) {
			score += 4
			if score >= 100 {
				return 100
			}
		}
	}
	if p.YearsExperience > 0 {
		score += min(10, p.YearsExperience)
	}
	if score > 100 {
		return 100
	}
	return score
}

func related(a, b string) bool {
	a, b = strings.ToLower(a), strings.ToLower(b)
	return strings.Contains(a, b) || strings.Contains(b, a)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
