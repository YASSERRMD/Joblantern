// Package safebrowsing submits confirmed-scam URLs to public anti-
// phishing programs (Google Safe Browsing, Microsoft SmartScreen,
// APWG). Each program has its own intake — we standardise the
// payload so the council can review once and ship to many.
package safebrowsing

import "time"

// Target is one downstream program.
type Target string

const (
	TargetGoogle    Target = "google-safe-browsing"
	TargetMicrosoft Target = "microsoft-smartscreen"
	TargetAPWG      Target = "apwg"
)

// Report is the abstract submission shape.
type Report struct {
	URL          string
	Category     string // "recruitment-scam"
	EvidenceURLs []string
	SubmittedAt  time.Time
	SubmitterOrg string
}

// CanonicalForm normalises the URL for downstream submission.
func (r *Report) CanonicalForm() {
	if r.SubmittedAt.IsZero() {
		r.SubmittedAt = time.Now().UTC()
	}
	if r.Category == "" {
		r.Category = "recruitment-scam"
	}
}
