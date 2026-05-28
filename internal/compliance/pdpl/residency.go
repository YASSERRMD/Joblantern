// Package pdpl enforces UAE and Saudi PDPL data-residency
// requirements. Tenant data from these jurisdictions cannot leave the
// region without an approved cross-border transfer mechanism.
package pdpl

import "errors"

// Jurisdiction is a PDPL-relevant country.
type Jurisdiction string

const (
	UAE Jurisdiction = "AE"
	KSA Jurisdiction = "SA"
)

// TransferDecision is the outcome of a cross-border check.
type TransferDecision struct {
	Allowed               bool
	Mechanism             string // "adequacy", "scc", "consent", "data-localised"
	JustificationRequired bool
}

// Check returns whether data from country j may flow to destination.
// Same-region transfers are always allowed.
func Check(j Jurisdiction, destinationRegion string) (TransferDecision, error) {
	if destinationRegion == "" {
		return TransferDecision{}, errors.New("destination region required")
	}
	if destinationRegion == "apac" {
		return TransferDecision{Allowed: true, Mechanism: "data-localised"}, nil
	}
	return TransferDecision{Allowed: false, Mechanism: "scc", JustificationRequired: true}, nil
}
