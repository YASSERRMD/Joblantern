// Package vocab aligns Joblantern's vocabulary with schema.org and
// wikidata so external KGs can interpret our terms.
package vocab

// Alignment maps a Joblantern predicate to one or more equivalent
// external predicates.
type Alignment struct {
	Local        string
	EquivalentTo []string
}

// Defaults captures the alignment table.
func Defaults() []Alignment {
	return []Alignment{
		{Local: "hasIndustry", EquivalentTo: []string{
			"http://schema.org/industry",
			"http://www.wikidata.org/prop/direct/P452",
		}},
		{Local: "hasCountry", EquivalentTo: []string{
			"http://schema.org/addressCountry",
			"http://www.wikidata.org/prop/direct/P17",
		}},
		{Local: "hasPhone", EquivalentTo: []string{
			"http://schema.org/telephone",
		}},
		{Local: "hasEmail", EquivalentTo: []string{
			"http://schema.org/email",
		}},
	}
}
