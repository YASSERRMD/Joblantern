package crosslink

import "testing"

func TestDetectCrossLinkedFixture(t *testing.T) {
	job := Job{
		Phones:  []string{"+97150-555-1234"},
		Emails:  []string{"hr@quickjobs-uae.com"},
		Domains: []string{"quickjobs-uae.com"},
	}
	rental := Rental{
		Phones:  []string{"+97150-555-1234"}, // same phone, different vertical
		Domains: []string{"luxe-villas-dubai.com"},
	}
	got := Detect(job, rental)
	if len(got) != 1 || got[0].Kind != "phone" {
		t.Fatalf("expected one phone match, got %#v", got)
	}
}

func TestDetectNoOverlap(t *testing.T) {
	if got := Detect(Job{Phones: []string{"a"}}, Rental{Phones: []string{"b"}}); len(got) != 0 {
		t.Fatalf("expected no matches, got %#v", got)
	}
}
