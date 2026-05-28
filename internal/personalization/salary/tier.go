// Package salary compares an offered salary against the realistic
// band for someone with the candidate's role and years of experience.
//
// The output is a personalised "is this salary plausible for *me*"
// adjustment, not a generic too-high / too-low alarm.
package salary

// Band is the inferred plausible salary range for a candidate at a
// destination country.
type Band struct {
	Country  string
	Currency string
	RoleP25  float64
	RoleP50  float64
	RoleP75  float64
}

// Bump computes a risk adjustment based on how far the offered
// salary deviates from the candidate-specific band.
//
//   - Below P25: reduces risk (legit jobs often underpay newcomers).
//   - At P25..P75: neutral.
//   - Above P75: increases risk (over-offer for the candidate level
//     is a classic scam pattern).
func Bump(offered float64, b Band) int {
	if offered <= 0 || b.RoleP50 == 0 {
		return 0
	}
	if offered < b.RoleP25 {
		return -5
	}
	if offered > b.RoleP75*2 {
		return +20
	}
	if offered > b.RoleP75 {
		return +8
	}
	return 0
}
