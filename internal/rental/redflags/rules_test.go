package redflags

import (
	"testing"
)

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
		{"western union deposit", Submission{DepositMethod: "Western Union only"}, "deposit-wire-only"},
		{"moneygram deposit", Submission{DepositMethod: "pay via MoneyGram"}, "deposit-wire-only"},
		{"gift card deposit", Submission{DepositMethod: "Amazon Gift Card"}, "deposit-wire-only"},
		{"deposit method case-insensitive", Submission{DepositMethod: "WESTERN UNION"}, "deposit-wire-only"},
		{"no viewing offered", Submission{NoInPersonViewingOffered: true}, "no-viewing"},
		{"reverse image hit", Submission{ReverseImageHits: 1}, "reverse-image-hit"},
		{"rent way below market", Submission{RentBelowMarketRatio: 0.4}, "rent-way-below-market"},
		{"urgency pressure", Submission{UrgencyPhrases: 3}, "urgency-pressure"},
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
		{"deposit method empty", Submission{DepositMethod: ""}, 0},
		{"benign deposit method", Submission{DepositMethod: "bank transfer"}, 0},
		{"reverse image zero hits", Submission{ReverseImageHits: 0}, 0},
		{"rent ratio zero is ignored", Submission{RentBelowMarketRatio: 0}, 0},
		{"rent ratio at 0.5 does not fire", Submission{RentBelowMarketRatio: 0.5}, 0},
		{"rent ratio just below 0.5 fires", Submission{RentBelowMarketRatio: 0.49}, 1},
		{"urgency at 2 does not fire", Submission{UrgencyPhrases: 2}, 0},
		{"urgency at 3 fires", Submission{UrgencyPhrases: 3}, 1},
		{"negative urgency does not fire", Submission{UrgencyPhrases: -1}, 0},
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
		DepositMethod:            "Western Union",
		NoInPersonViewingOffered: true,
		ReverseImageHits:         2,
		RentBelowMarketRatio:     0.4,
		UrgencyPhrases:           5,
	})

	want := map[string]int{
		"deposit-wire-only":     5,
		"no-viewing":            4,
		"reverse-image-hit":     4,
		"rent-way-below-market": 3,
		"urgency-pressure":      2,
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
		DepositMethod:            "bank transfer to verified account",
		NoInPersonViewingOffered: false,
		ReverseImageHits:         0,
		RentBelowMarketRatio:     0.9,
		UrgencyPhrases:           1,
	}); len(got) != 0 {
		t.Fatalf("expected no flags for clean submission, got %d: %+v", len(got), got)
	}
}
