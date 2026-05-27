// Package visa checks whether a claimed institution can plausibly
// issue the immigration documents the offer mentions.
package visa

// Eligibility describes the immigration document an institution can
// realistically issue.
type Eligibility struct {
	Country string
	Form    string // e.g. "I-20" (US), "CAS" (UK), "COE" (AU)
	Issued  bool   // true if the institution is recognised to issue it
}

// Check returns whether the claimed pathway is plausible for the
// claimed country and institution. A false here is a hard red flag.
func Check(country, form string, recognised []Eligibility) (Eligibility, bool) {
	for _, e := range recognised {
		if e.Country == country && e.Form == form {
			return e, e.Issued
		}
	}
	return Eligibility{}, false
}

// CanonicalForm maps a country to the standard form name.
func CanonicalForm(country string) string {
	switch country {
	case "US":
		return "I-20"
	case "UK":
		return "CAS"
	case "AU":
		return "COE"
	case "CA":
		return "LOA"
	case "DE":
		return "Zulassungsbescheid"
	}
	return ""
}
