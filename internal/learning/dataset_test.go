package learning_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/learning"
)

func TestExportCSV(t *testing.T) {
	rows := []learning.Labelled{
		{
			VerificationID: "v1",
			CompletedAt:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Country:        "AE", Risk: "red",
			PatternCodes: []string{"upfront_fee", "urgency"},
			Outcome:      "confirmed_scam",
		},
	}
	var buf bytes.Buffer
	if err := learning.ExportCSV(&buf, rows); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if !strings.Contains(s, "upfront_fee,urgency") || !strings.Contains(s, "confirmed_scam") {
		t.Fatalf("missing fields: %s", s)
	}
}

func TestEffectiveness(t *testing.T) {
	rows := []learning.Labelled{}
	for i := 0; i < 25; i++ {
		rows = append(rows, learning.Labelled{
			PatternCodes: []string{"upfront_fee"},
			Outcome:      "confirmed_scam",
		})
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, learning.Labelled{
			PatternCodes: []string{"upfront_fee"},
			Outcome:      "confirmed_legit",
		})
	}
	// noisy rule with mostly legit
	for i := 0; i < 30; i++ {
		rows = append(rows, learning.Labelled{
			PatternCodes: []string{"no_experience_needed"},
			Outcome:      "confirmed_legit",
		})
	}
	for i := 0; i < 5; i++ {
		rows = append(rows, learning.Labelled{
			PatternCodes: []string{"no_experience_needed"},
			Outcome:      "confirmed_scam",
		})
	}
	scores := learning.Effectiveness(rows)
	byCode := map[string]learning.EffectivenessScore{}
	for _, s := range scores {
		byCode[s.Code] = s
	}
	up := byCode["upfront_fee"]
	if up.Precision < 0.6 || up.Recommendation != "keep" {
		t.Errorf("upfront_fee score wrong: %+v", up)
	}
	noex := byCode["no_experience_needed"]
	if noex.Precision > 0.4 || noex.Recommendation != "consider_removal" {
		t.Errorf("no_experience_needed score wrong: %+v", noex)
	}
}
