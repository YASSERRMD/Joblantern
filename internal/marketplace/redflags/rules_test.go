package redflags

import "testing"

func hasCode(flags []Flag, code string) bool {
	for _, f := range flags {
		if f.Code == code {
			return true
		}
	}
	return false
}

// TestScanIndividualRules exercises each rule in isolation, asserting both the
// code that should fire and that no unrelated rule trips on the same input.
func TestScanIndividualRules(t *testing.T) {
	tests := []struct {
		name     string
		sub      Submission
		wantCode string
	}{
		{"western union payment", Submission{PaymentMethod: "Western Union only"}, "non-escrow-payment"},
		{"gift card payment", Submission{PaymentMethod: "Steam Gift Card"}, "non-escrow-payment"},
		{"crypto only payment", Submission{PaymentMethod: "Crypto only please"}, "non-escrow-payment"},
		{"payment method case-insensitive", Submission{PaymentMethod: "WESTERN UNION"}, "non-escrow-payment"},
		{"shipping via seller", Submission{BuyerPaysShippingViaSeller: true}, "shipping-via-seller"},
		{"advance fee requested", Submission{AdvanceFeeRequested: true}, "advance-fee"},
		{"price way below market", Submission{PriceBelowMedianRatio: 0.2}, "price-way-below-market"},
		{"off-platform pressure", Submission{OfferOnlyByDirect: true}, "off-platform-pressure"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := Scan(tt.sub)
			if !hasCode(flags, tt.wantCode) {
				t.Fatalf("expected flag %q, got %+v", tt.wantCode, flags)
			}
			if len(flags) != 1 {
				t.Fatalf("expected exactly 1 flag for isolated input, got %d: %+v", len(flags), flags)
			}
		})
	}
}

// TestScanBoundaries pins the exact thresholds so off-by-one regressions are caught.
func TestScanBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		sub     Submission
		wantLen int
	}{
		{"payment method empty", Submission{PaymentMethod: ""}, 0},
		{"benign payment method", Submission{PaymentMethod: "platform escrow"}, 0},
		{"price ratio zero is ignored", Submission{PriceBelowMedianRatio: 0}, 0},
		{"price ratio at 0.3 does not fire", Submission{PriceBelowMedianRatio: 0.3}, 0},
		{"price ratio just below 0.3 fires", Submission{PriceBelowMedianRatio: 0.29}, 1},
		{"negative price ratio does not fire", Submission{PriceBelowMedianRatio: -1}, 0},
		{"shipping method alone is inert", Submission{ShippingMethod: "courier"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Scan(tt.sub); len(got) != tt.wantLen {
				t.Fatalf("expected %d flags, got %d: %+v", tt.wantLen, len(got), got)
			}
		})
	}
}

// TestScanSeverities locks the severity assigned to each code.
func TestScanSeverities(t *testing.T) {
	flags := Scan(Submission{
		PaymentMethod:              "Western Union only",
		BuyerPaysShippingViaSeller: true,
		AdvanceFeeRequested:        true,
		PriceBelowMedianRatio:      0.2,
		OfferOnlyByDirect:          true,
	})

	want := map[string]int{
		"non-escrow-payment":     5,
		"shipping-via-seller":    4,
		"advance-fee":            5,
		"price-way-below-market": 3,
		"off-platform-pressure":  2,
	}

	if len(flags) != len(want) {
		t.Fatalf("expected %d flags, got %d: %+v", len(want), len(flags), flags)
	}
	for _, f := range flags {
		sev, ok := want[f.Code]
		if !ok {
			t.Fatalf("unexpected flag code %q", f.Code)
		}
		if f.Severity != sev {
			t.Errorf("code %q: expected severity %d, got %d", f.Code, sev, f.Severity)
		}
		if f.Message == "" {
			t.Errorf("code %q: expected non-empty message", f.Code)
		}
	}
}

// TestScanCleanSubmission confirms a legitimate listing produces no flags.
func TestScanCleanSubmission(t *testing.T) {
	if got := Scan(Submission{
		PaymentMethod:         "platform escrow",
		ShippingMethod:        "tracked courier",
		PriceBelowMedianRatio: 0.9,
	}); len(got) != 0 {
		t.Fatalf("expected no flags for clean submission, got %d: %+v", len(got), got)
	}
}
