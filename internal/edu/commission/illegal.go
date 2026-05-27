// Package commission flags illegal recruitment commissions for
// education agents. Many origin-country regulators (e.g., the
// Philippines POEA, India's MEA) cap or prohibit agent commissions
// charged to the student.
package commission

// Rule captures the per-country position on student-paid commissions.
type Rule struct {
	Country         string
	StudentMaxFeeUSD float64
	Notes           string
}

// Defaults is a seed table.
var Defaults = []Rule{
	{Country: "PH", StudentMaxFeeUSD: 0, Notes: "Education agent commissions charged to the student are prohibited."},
	{Country: "IN", StudentMaxFeeUSD: 0, Notes: "Agent fees beyond documented service charges are flagged."},
	{Country: "NP", StudentMaxFeeUSD: 0, Notes: "Education consultancies cannot charge for placement."},
	{Country: "BD", StudentMaxFeeUSD: 0, Notes: "Same restriction as IN/NP."},
}

// Flag is true when the alleged commission exceeds the per-country cap.
func Flag(country string, allegedUSD float64) (bool, string) {
	for _, r := range Defaults {
		if r.Country == country && allegedUSD > r.StudentMaxFeeUSD {
			return true, r.Notes
		}
	}
	return false, ""
}
