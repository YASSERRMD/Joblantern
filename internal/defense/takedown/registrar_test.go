package takedown

import (
	"strings"
	"testing"
	"time"
)

func TestPacketContainsAllSections(t *testing.T) {
	p := Packet{
		Registrar:    "namecheap.com",
		Domain:       "quickjobs-uae.example",
		ReportedAt:   time.Date(2026, 5, 27, 0, 0, 0, 0, time.UTC),
		ContactEmail: "abuse@joblantern.org",
		ContactOrg:   "Joblantern",
		Evidence: []Evidence{
			{Kind: "verdict", Description: "Red verdict, confidence 0.94", URL: "https://joblantern.org/v/123", SHA256: "deadbeef"},
		},
	}
	out, err := Render(p)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"Abuse report", "namecheap", "quickjobs-uae.example", "deadbeef"} {
		if !strings.Contains(out, s) {
			t.Errorf("packet missing %q", s)
		}
	}
}
