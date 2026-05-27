// Package billing is the scaffold for tier-based pricing.
// Low-income NGOs run on the free tier; commercial recruiters pay.
//
// Joblantern does not bill directly in this scaffold — production
// wires up Stripe, but the policy layer (eligibility, caps) lives
// here so it can be reasoned about independently of provider.
package billing

// Plan is one subscription plan.
type Plan struct {
	ID                string
	DisplayName       string
	MonthlyPriceUSD   float64
	IncludedVerdicts  int
	OverageCentsPer1k int
	NGODiscount       bool
}

// Plans is the canonical list.
var Plans = []Plan{
	{ID: "ngo-free", DisplayName: "NGO Free", IncludedVerdicts: 1000, NGODiscount: true},
	{ID: "pro", DisplayName: "Pro", MonthlyPriceUSD: 49, IncludedVerdicts: 25000, OverageCentsPer1k: 200},
	{ID: "custom", DisplayName: "Custom", IncludedVerdicts: -1},
}

// Eligible returns true if the tenant qualifies for the plan.
func Eligible(tier string, planID string) bool {
	switch planID {
	case "ngo-free":
		return tier == "free"
	case "pro":
		return tier == "pro"
	case "custom":
		return tier == "custom"
	}
	return false
}
