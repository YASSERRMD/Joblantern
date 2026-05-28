// Package redflags scores rental-specific scam patterns.
package redflags

import (
	"strings"
)

// Flag is a structured red-flag finding.
type Flag struct {
	Code     string
	Severity int // 1 (info) .. 5 (critical)
	Message  string
}

// Scan runs the rental rule pack against a submission.
func Scan(s Submission) []Flag {
	var flags []Flag

	if s.DepositMethod != "" {
		dm := strings.ToLower(s.DepositMethod)
		if strings.Contains(dm, "western union") || strings.Contains(dm, "moneygram") || strings.Contains(dm, "gift card") {
			flags = append(flags, Flag{"deposit-wire-only", 5, "Deposit by Western Union / MoneyGram / gift card is a near-universal scam signal."})
		}
	}
	if s.NoInPersonViewingOffered {
		flags = append(flags, Flag{"no-viewing", 4, "Landlord refuses any in-person or video viewing before deposit."})
	}
	if s.ReverseImageHits > 0 {
		flags = append(flags, Flag{"reverse-image-hit", 4, "Listing photos appear on other listings under different landlords."})
	}
	if s.RentBelowMarketRatio > 0 && s.RentBelowMarketRatio < 0.5 {
		flags = append(flags, Flag{"rent-way-below-market", 3, "Rent is less than half the local median for similar units."})
	}
	if s.UrgencyPhrases > 2 {
		flags = append(flags, Flag{"urgency-pressure", 2, "Listing leans heavily on urgency phrasing (\"must decide today\", \"others waiting\")."})
	}
	return flags
}

// Submission is the rental input the rule pack consumes.
type Submission struct {
	DepositMethod            string
	NoInPersonViewingOffered bool
	ReverseImageHits         int
	RentBelowMarketRatio     float64
	UrgencyPhrases           int
}
