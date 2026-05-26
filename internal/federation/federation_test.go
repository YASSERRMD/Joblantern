package federation_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/federation"
)

func TestSignVerify_RoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sig := federation.Signal{
		OriginURL:    "https://acme.ngo/",
		CompanyName:  "Phantom Recruiters",
		Country:      "AE",
		PatternCodes: []string{"upfront_fee", "identity_documents_upfront"},
		Verdict:      "red",
		IssuedAt:     time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
	}
	env, err := federation.Sign(priv, sig)
	if err != nil {
		t.Fatal(err)
	}
	if err := federation.Verify(pub, env); err != nil {
		t.Fatalf("verify: %v", err)
	}
}

func TestVerify_Tamper(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	env, _ := federation.Sign(priv, federation.Signal{
		CompanyName: "Acme", Country: "AE", Verdict: "red", IssuedAt: time.Now().UTC(),
	})
	env.Signal.CompanyName = "Acme Holdings"
	if err := federation.Verify(pub, env); err == nil {
		t.Fatal("expected verify to fail after tamper")
	}
}

func TestFingerprint_Stable(t *testing.T) {
	a := federation.Signal{
		CompanyName:  "Acme",
		Country:      "AE",
		PatternCodes: []string{"a", "b"},
	}
	b := federation.Signal{
		CompanyName:  "Acme",
		Country:      "AE",
		PatternCodes: []string{"b", "a"},
	}
	if a.Fingerprint() != b.Fingerprint() {
		t.Errorf("fingerprint changed with code ordering")
	}
}
