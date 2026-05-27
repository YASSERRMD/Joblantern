// Package generator produces realistic synthetic scam listings to
// exercise the agent stack. The fixtures are not adversarial — they
// are average-case scams the agent should catch every time.
package generator

import "time"

// Listing is a generated test listing.
type Listing struct {
	ID        string
	Title     string
	Body      string
	Country   string
	Industry  string
	Phone     string
	Email     string
	Domain    string
	Salary    float64
	Currency  string
	GroundTruth string // "scam", "ambiguous", "legit"
}

// Defaults returns the seed list. Production generators expand from
// templates per origin country and industry.
func Defaults(now time.Time) []Listing {
	return []Listing{
		{ID: "syn-1", Title: "URGENT Driver wanted Dubai", Body: "Pay AED 18000/mo. Apply via WhatsApp now.", Country: "AE", Industry: "transport", Phone: "+971-55-555-1111", GroundTruth: "scam"},
		{ID: "syn-2", Title: "Nurse, KSA — bring passport copy", Body: "No interview needed. Visa on arrival.", Country: "SA", Industry: "health", GroundTruth: "scam"},
		{ID: "syn-3", Title: "Software Engineer, Doha", Body: "Standard offer letter attached. EU recruiter.", Country: "QA", Industry: "tech", GroundTruth: "legit"},
	}
}
