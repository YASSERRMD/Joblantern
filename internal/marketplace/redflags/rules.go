// Package redflags scores marketplace-specific scam patterns.
package redflags

import "strings"

// Flag is a finding.
type Flag struct {
	Code     string
	Severity int
	Message  string
}

// Submission is the input the rule pack consumes.
type Submission struct {
	PaymentMethod              string
	ShippingMethod             string
	PriceBelowMedianRatio      float64
	OfferOnlyByDirect          bool
	BuyerPaysShippingViaSeller bool
	AdvanceFeeRequested        bool
}

// Scan runs the rule pack.
func Scan(s Submission) []Flag {
	var flags []Flag
	if pm := strings.ToLower(s.PaymentMethod); strings.Contains(pm, "western union") || strings.Contains(pm, "gift card") || strings.Contains(pm, "crypto only") {
		flags = append(flags, Flag{"non-escrow-payment", 5, "Payment outside platform escrow is a near-universal marketplace scam pattern."})
	}
	if s.BuyerPaysShippingViaSeller {
		flags = append(flags, Flag{"shipping-via-seller", 4, "Buyer is asked to pay a 'shipping fee' via the seller — typical shipping-diversion."})
	}
	if s.AdvanceFeeRequested {
		flags = append(flags, Flag{"advance-fee", 5, "Seller asks for an upfront non-refundable fee."})
	}
	if s.PriceBelowMedianRatio > 0 && s.PriceBelowMedianRatio < 0.3 {
		flags = append(flags, Flag{"price-way-below-market", 3, "Asking price is less than 30% of median for similar items."})
	}
	if s.OfferOnlyByDirect {
		flags = append(flags, Flag{"off-platform-pressure", 2, "Seller refuses to communicate on the platform."})
	}
	return flags
}
