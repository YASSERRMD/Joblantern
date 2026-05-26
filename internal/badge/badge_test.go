package badge_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/badge"
)

func TestEncodeDecode(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	c := badge.Claims{
		BadgeID: "b-1", OrgID: "o-1", OrgName: "Acme",
		VerificationID: "v-1", Risk: "green",
		IssuedAt:   time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(24 * time.Hour),
		Issuer:     "https://joblantern.example/",
		TrustLevel: "observed",
	}
	tok, err := badge.Encode(priv, c)
	if err != nil {
		t.Fatal(err)
	}
	got, err := badge.Decode(pub, tok)
	if err != nil {
		t.Fatal(err)
	}
	if got.BadgeID != "b-1" || got.Risk != "green" {
		t.Errorf("%+v", got)
	}
}

func TestDecode_Expired(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	c := badge.Claims{
		BadgeID: "b-1", OrgID: "o-1", OrgName: "Acme",
		VerificationID: "v-1", Risk: "green",
		IssuedAt:   time.Now().UTC().Add(-2 * time.Hour),
		ExpiresAt:  time.Now().UTC().Add(-1 * time.Hour),
		Issuer:     "https://joblantern.example/",
		TrustLevel: "observed",
	}
	tok, _ := badge.Encode(priv, c)
	_, err := badge.Decode(pub, tok)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected expired error, got %v", err)
	}
}

func TestDecode_Tamper(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	c := badge.Claims{
		BadgeID: "b-1", OrgID: "o-1", VerificationID: "v-1",
		Risk: "green", IssuedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Issuer:    "https://j/",
	}
	tok, _ := badge.Encode(priv, c)
	// Flip one base64 char in the payload portion.
	idx := strings.Index(tok, ".")
	tampered := tok[:5] + "A" + tok[6:idx] + tok[idx:]
	if _, err := badge.Decode(pub, tampered); err == nil {
		t.Fatal("expected signature mismatch")
	}
}
