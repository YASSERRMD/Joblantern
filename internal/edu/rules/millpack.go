// Package rules contains the diploma-mill rule pack.
package rules

import "strings"

// Flag is a finding.
type Flag struct {
	Code     string
	Severity int
	Message  string
}

// Submission is the input.
type Submission struct {
	Institution           string
	Program               string
	OffersLifeExperience  bool
	AccreditorClaim       string
	UnaccreditedConfirmed bool
	WebsiteAgeDays        int
	UncorroboratedProgram bool
}

// Scan applies the rule pack.
func Scan(s Submission) []Flag {
	var flags []Flag
	if s.UnaccreditedConfirmed {
		flags = append(flags, Flag{"unaccredited", 5, "Institution does not appear in any recognised national accreditation registry."})
	}
	if s.OffersLifeExperience {
		flags = append(flags, Flag{"life-experience-degree", 5, "Institution offers degrees for \"life experience\" — a near-universal diploma mill marker."})
	}
	if low := strings.ToLower(s.AccreditorClaim); strings.Contains(low, "accreditation council for online") || strings.Contains(low, "universal accrediting") {
		flags = append(flags, Flag{"dubious-accreditor", 4, "Accreditor name matches known mill-accreditor pattern."})
	}
	if s.WebsiteAgeDays > 0 && s.WebsiteAgeDays < 365 && s.UnaccreditedConfirmed {
		flags = append(flags, Flag{"young-domain-unaccredited", 3, "Institution website registered <12 months ago and not accredited."})
	}
	if s.UncorroboratedProgram {
		flags = append(flags, Flag{"program-not-in-catalogues", 2, "Claimed program not found in published course catalogues."})
	}
	return flags
}
