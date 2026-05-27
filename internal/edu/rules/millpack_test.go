package rules

import "testing"

func TestPacificWesternUniversityFixture(t *testing.T) {
	s := Submission{
		Institution:           "Pacific Western University",
		Program:               "Doctorate of Management",
		OffersLifeExperience:  true,
		AccreditorClaim:       "Accreditation Council for Online Academia",
		UnaccreditedConfirmed: true,
		WebsiteAgeDays:        180,
		UncorroboratedProgram: true,
	}
	flags := Scan(s)
	if len(flags) < 4 {
		t.Fatalf("expected >=4 flags for known mill, got %d", len(flags))
	}
}

func TestAccreditedNoFlags(t *testing.T) {
	flags := Scan(Submission{Institution: "ETH Zürich"})
	if len(flags) != 0 {
		t.Fatalf("legitimate institution should produce no flags, got %d", len(flags))
	}
}
