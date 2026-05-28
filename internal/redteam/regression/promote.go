// Package regression promotes a missed adversarial fixture into a
// permanent regression test.
package regression

import "time"

// Promotion is one promoted fixture entry.
type Promotion struct {
	FixtureID  string
	BodySHA256 string
	Mutators   []string
	PromotedAt time.Time
	PromotedBy string
	Rationale  string
}

// File is the on-disk format we serialise to internal/redteam/corpus.
type File struct {
	Promotions []Promotion
}

// Add appends a promotion entry deduplicated by FixtureID.
func (f *File) Add(p Promotion) {
	for i, e := range f.Promotions {
		if e.FixtureID == p.FixtureID {
			f.Promotions[i] = p
			return
		}
	}
	f.Promotions = append(f.Promotions, p)
}
