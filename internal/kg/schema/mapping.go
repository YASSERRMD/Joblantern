// Package schema maps verifications and evidence into the triple
// store schema. The IRIs use a stable prefix that we publish in our
// vocabulary document so external KGs can link to us.
package schema

import "fmt"

// Base is the canonical IRI prefix.
const Base = "https://joblantern.org/kg/"

// Predicate is a typed predicate IRI helper.
type Predicate string

const (
	PredHasRiskBand Predicate = "hasRiskBand"
	PredHasIndustry Predicate = "hasIndustry"
	PredHasCountry  Predicate = "hasCountry"
	PredHasPhone    Predicate = "hasPhone"
	PredHasEmail    Predicate = "hasEmail"
	PredHasDirector Predicate = "hasDirector"
	PredHasAddress  Predicate = "hasAddress"
	PredHasEvidence Predicate = "hasEvidence"
	PredVerifiedAt  Predicate = "verifiedAt"
)

// VerdictIRI returns the IRI for a verdict id.
func VerdictIRI(id string) string { return fmt.Sprintf("%sverdict/%s", Base, id) }

// PredIRI returns the IRI for a predicate.
func PredIRI(p Predicate) string { return Base + "predicate/" + string(p) }

// EntityIRI returns the IRI for an entity node.
func EntityIRI(kind, id string) string { return fmt.Sprintf("%sentity/%s/%s", Base, kind, id) }
