// Package combined produces a single user-facing verdict when a
// submission includes both a job offer and a rental listing.
package combined

// Combined is the joined output.
type Combined struct {
	JobBand        string
	RentalBand     string
	OverallBand    string
	CrossLinkFlags []string
	Summary        string
}

// Merge picks the worst of the two bands, then escalates one notch
// if any cross-link flags were detected.
func Merge(job, rental string, crossLink []string) Combined {
	overall := worst(job, rental)
	if len(crossLink) > 0 {
		overall = bumpUp(overall)
	}
	return Combined{
		JobBand:        job,
		RentalBand:     rental,
		OverallBand:    overall,
		CrossLinkFlags: crossLink,
	}
}

func worst(a, b string) string {
	r := map[string]int{"green": 0, "yellow": 1, "red": 2}
	if r[a] >= r[b] {
		return a
	}
	return b
}

func bumpUp(b string) string {
	switch b {
	case "green":
		return "yellow"
	case "yellow":
		return "red"
	}
	return b
}
