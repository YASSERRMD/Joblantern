package fit

import "testing"

func TestExactRoleMatchScores60Plus(t *testing.T) {
	got := Score(Listing{Role: "welder", Description: "MIG and TIG welding"}, Profile{Role: "Welder", YearsExperience: 5, Skills: []string{"MIG"}})
	if got < 60 {
		t.Fatalf("expected >= 60, got %d", got)
	}
}

func TestUnrelatedRole(t *testing.T) {
	got := Score(Listing{Role: "Welder"}, Profile{Role: "Software engineer"})
	if got > 50 {
		t.Fatalf("expected unrelated to be <=50, got %d", got)
	}
}
