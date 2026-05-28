package sandbox

import (
	"testing"
	"time"
)

func TestProvisionAutomation(t *testing.T) {
	now := time.Now()
	s := Spec{PartnerID: "uni-1", Region: "eu", StartAt: now, EndAt: now.Add(180 * 24 * time.Hour), Tier: "C"}
	if err := s.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestRejectInvalidWindow(t *testing.T) {
	now := time.Now()
	// Same instant for start and end == a zero-length window, which
	// Validate must reject (end is not strictly after start).
	if err := (Spec{PartnerID: "uni-1", StartAt: now, EndAt: now, Tier: "B"}).Validate(); err == nil {
		t.Fatal("expected error for zero-length window")
	}
}
