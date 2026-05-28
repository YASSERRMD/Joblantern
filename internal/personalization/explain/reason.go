// Package explain renders the "this offer is unusual for your
// profile because..." sentence that accompanies a personalized
// verdict. Phrasing is deliberately concrete and non-judgmental.
package explain

import (
	"fmt"
	"strings"
)

// Factors is the structured input to the explainer.
type Factors struct {
	RoleMismatch          bool
	SalaryFarAboveBand    bool
	DestinationAnomaly    bool
	YearsBelowRequirement bool
}

// Render produces a one-paragraph explanation. Empty factors yield an
// empty string so the UI can skip the section entirely.
func Render(f Factors) string {
	var bits []string
	if f.RoleMismatch {
		bits = append(bits, "the listing's role does not match the one in your CV")
	}
	if f.SalaryFarAboveBand {
		bits = append(bits, "the salary is far above the typical band for your level")
	}
	if f.DestinationAnomaly {
		bits = append(bits, "this destination country is unusual for someone with your background")
	}
	if f.YearsBelowRequirement {
		bits = append(bits, "the offer expects more years of experience than your CV lists")
	}
	if len(bits) == 0 {
		return ""
	}
	return fmt.Sprintf("This offer is unusual for your profile because %s.", strings.Join(bits, "; "))
}
