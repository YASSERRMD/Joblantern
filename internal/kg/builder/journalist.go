// Package builder is the simple, journalist-friendly SPARQL query
// builder UI. It exposes a handful of pre-built starter queries and
// lets the user fill in country/industry/date filters.
package builder

import "strings"

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
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}
