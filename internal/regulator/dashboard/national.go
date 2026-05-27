// Package dashboard renders the regulator-facing dashboard. By
// policy it exposes only national-level aggregates and never raw
// per-applicant submissions.
package dashboard

import "time"

// Stats is the canonical national view served to a regulator account.
type Stats struct {
	Country               string    `json:"country"`
	GeneratedAt           time.Time `json:"generated_at"`
	TotalVerdicts         int       `json:"total_verdicts"`
	RedVerdicts           int       `json:"red_verdicts"`
	YellowVerdicts        int       `json:"yellow_verdicts"`
	GreenVerdicts         int       `json:"green_verdicts"`
	TopRedFlags           []Bucket  `json:"top_red_flags"`
	TopDestinationIndustries []Bucket `json:"top_destination_industries"`
}

// Bucket is a label/count pair.
type Bucket struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// MinBucketSize prevents micro-aggregate disclosure that could
// re-identify individuals. Any bucket below this floor is coalesced
// into "other".
const MinBucketSize = 5

// Coalesce removes buckets below MinBucketSize, summing their counts
// into a single "other" entry.
func Coalesce(b []Bucket) []Bucket {
	out := make([]Bucket, 0, len(b)+1)
	other := 0
	for _, x := range b {
		if x.Count < MinBucketSize {
			other += x.Count
			continue
		}
		out = append(out, x)
	}
	if other > 0 {
		out = append(out, Bucket{Label: "other", Count: other})
	}
	return out
}
