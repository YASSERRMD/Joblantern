package redflags

import "testing"

func TestUsedPhoneFixtureCatchesNetwork(t *testing.T) {
	flags := Scan(Submission{
		PaymentMethod:              "Western Union only",
		BuyerPaysShippingViaSeller: true,
		AdvanceFeeRequested:        true,
		PriceBelowMedianRatio:      0.2,
		OfferOnlyByDirect:          true,
	})
	if len(flags) < 4 {
		t.Fatalf("expected >=4 flags, got %d", len(flags))
	}
}

func TestLegitimatePassesClean(t *testing.T) {
	if got := Scan(Submission{PaymentMethod: "platform escrow"}); len(got) != 0 {
		t.Fatalf("expected no flags for clean fixture, got %d", len(got))
	}
}
