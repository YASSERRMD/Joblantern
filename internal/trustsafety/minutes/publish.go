// Package minutes publishes council minutes where appropriate. The
// charter allows redaction for personal-data protection; the policy
// is "redact, do not omit".
package minutes

import "time"

// Minute is one published meeting record.
type Minute struct {
	ID        string
	MeetingAt time.Time
	Quorum    bool
	Topics    []Topic
	PublishedAt time.Time
}

// Topic is one discussion item.
type Topic struct {
	Title   string
	Outcome string
	CasesRef []string
	Redacted bool
}

// Publish marks the minute as published.
func Publish(m *Minute, at time.Time) {
	m.PublishedAt = at.UTC()
}
