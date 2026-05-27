// Package bias is the audit harness for personalized verdicts.
//
// We must not produce systematically harsher verdicts for protected
// categories. The audit runs across a battery of profile pairs that
// differ only in protected attributes and asserts that the resulting
// verdict bands match.
package bias

// Pair is a controlled comparison: the same listing assessed against
// two profiles that differ only in one protected attribute.
type Pair struct {
	Listing string
	A, B    Profile
}

// Profile is the minimum attribute set for the audit harness.
type Profile struct {
	Role            string
	YearsExperience int
	OriginCountry   string
	Gender          string
	AgeBracket      string
}

// Result records the audit outcome for one pair.
type Result struct {
	Pair      Pair
	BandA     string
	BandB     string
	Equivalent bool
}

// Equivalent returns true when the verdict bands match across the
// pair. Production code returns the divergent set so the Trust &
// Safety council can review.
func (r Result) IsEquivalent() bool { return r.BandA == r.BandB }
