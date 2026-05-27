// Package builder is the simple, journalist-friendly SPARQL query
// builder UI. It exposes a handful of pre-built starter queries and
// lets the user fill in country/industry/date filters.
package builder

import "fmt"

// Starter is one pre-built query.
type Starter struct {
	ID     string
	Title  string
	Query  string // SPARQL with {placeholders}
	Params []string
}

// Starters returns the canonical list.
func Starters() []Starter {
	return []Starter{
		{
			ID:    "country-confirmed",
			Title: "Confirmed scam companies in {COUNTRY}",
			Query: `PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT ?company ?industry WHERE {
  ?company jl:hasCountry "{COUNTRY}" ;
           jl:hasRiskBand "red" ;
           jl:hasIndustry ?industry .
} LIMIT 200`,
			Params: []string{"COUNTRY"},
		},
		{
			ID:    "shared-phone-across-countries",
			Title: "Companies in {COUNTRY_A} sharing a phone with any company in {COUNTRY_B}",
			Query: `PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT ?a ?b ?phone WHERE {
  ?a jl:hasCountry "{COUNTRY_A}" ; jl:hasPhone ?phone .
  ?b jl:hasCountry "{COUNTRY_B}" ; jl:hasPhone ?phone .
  FILTER (?a != ?b)
} LIMIT 500`,
			Params: []string{"COUNTRY_A", "COUNTRY_B"},
		},
	}
}

// Fill substitutes the placeholders in a starter.
func Fill(s Starter, args map[string]string) string {
	out := s.Query
	for k, v := range args {
		out = replace(out, "{"+k+"}", v)
	}
	return out
}

func replace(s, old, new string) string {
	if old == "" {
		return s
	}
	out := ""
	for {
		i := indexOf(s, old)
		if i < 0 {
			return out + s
		}
		out += s[:i] + new
		s = s[i+len(old):]
	}
}

func indexOf(s, sub string) int {
	if len(sub) == 0 {
		return 0
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func ensureFmt() string { return fmt.Sprintf("%s", "") }
