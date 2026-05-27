// Package coordinated integrates Joblantern with cooperatively-run
// takedown groups (StopNCII-style). Joblantern is a participant, not
// a member-of-record — we contribute and consume per-group rules.
package coordinated

// Group is one cooperative takedown group we participate in.
type Group struct {
	ID           string
	Name         string
	Contact      string
	URL          string
	MembershipID string
}

// Default returns the seed list of relevant groups, kept short and
// explicit so contracts and MOUs are visible per group.
func Default() []Group {
	return []Group{
		{ID: "apwg", Name: "Anti-Phishing Working Group", URL: "https://apwg.org"},
		{ID: "stopncii", Name: "StopNCII / SWGfL", URL: "https://stopncii.org"},
	}
}
