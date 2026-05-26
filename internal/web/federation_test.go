package web_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/federation"
	"github.com/yasserrmd/joblantern/internal/web"
)

func TestFederation_IngestRoundTrip(t *testing.T) {
	// Peer A signs, instance B verifies.
	pubA, privA, _ := ed25519.GenerateKey(rand.Reader)
	pubB, _, _ := ed25519.GenerateKey(rand.Reader)

	r := chi.NewRouter()
	lookup := func(url string) (ed25519.PublicKey, bool) {
		if url == "https://peer-a.example/" {
			return pubA, true
		}
		return nil, false
	}
	web.NewFederation(r, "instance-b", "https://b.example/", "0.1.0", pubB, lookup)
	srv := httptest.NewServer(r)
	defer srv.Close()

	// Build + sign a signal from peer A.
	env, err := federation.Sign(privA, federation.Signal{
		OriginURL:    "https://peer-a.example/",
		CompanyName:  "Phantom Recruiters",
		Country:      "AE",
		PatternCodes: []string{"upfront_fee"},
		Verdict:      "red",
		IssuedAt:     time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(env)
	resp, err := http.Post(srv.URL+"/a2a/signal", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("ingest got %d", resp.StatusCode)
	}

	// Recent endpoint should return our signal.
	rec, _ := http.Get(srv.URL + "/a2a/recent")
	var got []federation.Signal
	_ = json.NewDecoder(rec.Body).Decode(&got)
	_ = rec.Body.Close()
	if len(got) != 1 || got[0].CompanyName != "Phantom Recruiters" {
		t.Fatalf("got %+v", got)
	}
}

func TestFederation_UnknownPeerRejected(t *testing.T) {
	pubSelf, _, _ := ed25519.GenerateKey(rand.Reader)
	_, privUnknown, _ := ed25519.GenerateKey(rand.Reader)

	r := chi.NewRouter()
	web.NewFederation(r, "b", "https://b/", "0.1.0", pubSelf,
		func(string) (ed25519.PublicKey, bool) { return nil, false })
	srv := httptest.NewServer(r)
	defer srv.Close()

	env, _ := federation.Sign(privUnknown, federation.Signal{
		OriginURL: "https://stranger.example/", Country: "AE",
		Verdict: "red", IssuedAt: time.Now().UTC(),
	})
	body, _ := json.Marshal(env)
	resp, err := http.Post(srv.URL+"/a2a/signal", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}
