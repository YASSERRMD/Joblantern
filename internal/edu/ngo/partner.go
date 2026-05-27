// Package ngo provides the student-association partner surface so a
// university student union or migrant-support NGO can host a branded
// instance of the edu module for their members.
package ngo

import "time"

// Partner is a student association running an edu instance.
type Partner struct {
	ID         string
	Name       string
	Country    string
	Logo       string
	Theme      string
	Members    int
	ActivatedAt time.Time
}

// Welcome composes the partner-facing welcome string. Kept short so it
// fits on a single membership card or QR landing page.
func (p Partner) Welcome() string {
	if p.Name == "" {
		return "Welcome to Joblantern Edu"
	}
	return p.Name + " × Joblantern Edu"
}
