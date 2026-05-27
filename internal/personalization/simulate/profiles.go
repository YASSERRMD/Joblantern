// Package simulate lets a user explore "what would my verdict be if
// I were a different person?" Useful for awareness training, NGO
// teaching, and journalist demos.
package simulate

// Profile is a synthetic persona used to re-score a listing.
type Profile struct {
	Label           string
	Role            string
	YearsExperience int
	OriginCountry   string
}

// Bank is the shipped set of demo personas.
var Bank = []Profile{
	{Label: "Junior welder, Bangladesh", Role: "welder", YearsExperience: 2, OriginCountry: "BD"},
	{Label: "Senior welder, Bangladesh", Role: "welder", YearsExperience: 15, OriginCountry: "BD"},
	{Label: "Software engineer, India", Role: "software engineer", YearsExperience: 4, OriginCountry: "IN"},
	{Label: "Nurse, Philippines", Role: "registered nurse", YearsExperience: 6, OriginCountry: "PH"},
	{Label: "Domestic worker, Nepal", Role: "domestic worker", YearsExperience: 3, OriginCountry: "NP"},
}
