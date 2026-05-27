package cv

// Profile is the structured form extracted from a CV. The parser runs
// against a small local model — by policy nothing in here ever leaves
// the trust boundary unless the user explicitly opts in to research
// donation.
type Profile struct {
	Role            string   `json:"role"`
	YearsExperience int      `json:"years_experience"`
	OriginCountry   string   `json:"origin_country"`
	Skills          []string `json:"skills"`
	SeniorityHint   string   `json:"seniority_hint"`
}

// Parser is implemented by the local model wrapper.
type Parser interface {
	Parse(raw []byte, contentType string) (Profile, error)
}

// Heuristic is a deterministic fallback parser used when the local
// model is unavailable. It keys off filename + content-type and lets
// downstream code degrade gracefully.
type Heuristic struct{}

// Parse returns an empty profile. The agent treats an empty profile
// as "no personalisation requested".
func (Heuristic) Parse(raw []byte, ct string) (Profile, error) { return Profile{}, nil }
