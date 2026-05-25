package pattern_test

import (
	"testing"

	"github.com/yasserrmd/joblantern/internal/pattern"
)

func TestDefaultPackAnalyse_Scam(t *testing.T) {
	rp, err := pattern.DefaultPack()
	if err != nil {
		t.Fatal(err)
	}
	scam := `Urgent hiring for warehouse worker in Dubai! No experience needed.
Earn USD 500 per day. Send passport copy and AED 2000 registration fee to WhatsApp +971500000000.
Government approved.`
	r := rp.Analyse(scam)
	if len(r.RedFlags) < 3 {
		t.Fatalf("expected >=3 red flags, got %d: %+v", len(r.RedFlags), r.RedFlags)
	}
	if r.CompositeScore < 0.7 {
		t.Errorf("composite score too low: %v", r.CompositeScore)
	}
}

func TestDefaultPackAnalyse_Clean(t *testing.T) {
	rp, _ := pattern.DefaultPack()
	clean := `Software engineer, Joblantern Ltd, Dubai. 4+ years Go. Apply via careers.joblantern.org.`
	r := rp.Analyse(clean)
	if len(r.RedFlags) > 1 {
		t.Fatalf("clean listing should not trigger many flags, got %+v", r.RedFlags)
	}
}

func TestLanguageMismatch(t *testing.T) {
	if mismatch, kind := pattern.LanguageMismatchCheck("Привет Дубай", "AE"); !mismatch || kind != "cyrillic" {
		t.Errorf("expected cyrillic mismatch, got %v %q", mismatch, kind)
	}
	if mismatch, _ := pattern.LanguageMismatchCheck("Hello Dubai", "AE"); mismatch {
		t.Error("expected no mismatch for english+AE")
	}
}
