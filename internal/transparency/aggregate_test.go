package transparency_test

import (
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/transparency"
)

func TestAggregate_DropsSmallCells(t *testing.T) {
	a := transparency.New()
	a.NoiseScale = 0 // disable noise for assertion stability
	day := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	in := make([]transparency.Verdict, 0, 12)
	// 6 AE/red — kept
	for i := 0; i < 6; i++ {
		in = append(in, transparency.Verdict{CompletedAt: day, Country: "AE", Risk: "red"})
	}
	// 2 PH/red — dropped (under MinCell=5)
	for i := 0; i < 2; i++ {
		in = append(in, transparency.Verdict{CompletedAt: day, Country: "PH", Risk: "red"})
	}
	rows := a.Aggregate(in)
	if len(rows) != 1 || rows[0].Country != "AE" || rows[0].Count != 6 {
		t.Fatalf("got %+v", rows)
	}
}

func TestAggregate_MissingCountryFallsToZZ(t *testing.T) {
	a := transparency.New()
	a.NoiseScale = 0
	day := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	in := make([]transparency.Verdict, 0, 5)
	for i := 0; i < 5; i++ {
		in = append(in, transparency.Verdict{CompletedAt: day, Risk: "yellow"})
	}
	rows := a.Aggregate(in)
	if len(rows) != 1 || rows[0].Country != "ZZ" {
		t.Fatalf("expected ZZ bucket, got %+v", rows)
	}
}
