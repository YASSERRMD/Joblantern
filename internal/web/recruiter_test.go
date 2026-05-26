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

	"github.com/yasserrmd/joblantern/internal/badge"
	"github.com/yasserrmd/joblantern/internal/web"
)

func TestBadge_IssueAndVerify(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	r := chi.NewRouter()
	web.NewBadgeIssuer(r, "https://joblantern.example/", priv, pub, 24*time.Hour)
	srv := httptest.NewServer(r)
	defer srv.Close()

	// Issue.
	body, _ := json.Marshal(map[string]string{
		"org_id":          "org-1",
		"org_name":        "Acme Recruiting",
		"verification_id": "v-1",
		"risk":            "green",
		"trust_level":     "vouched",
	})
	resp, err := http.Post(srv.URL+"/api/v1/recruiter/badges", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("issue status %d", resp.StatusCode)
	}
	var ir struct {
		Token  string       `json:"token"`
		Claims badge.Claims `json:"claims"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&ir)
	_ = resp.Body.Close()
	if ir.Token == "" || ir.Claims.Risk != "green" {
		t.Fatalf("got %+v", ir)
	}

	// Public verifier.
	resp2, err := http.Get(srv.URL + "/badge/" + ir.Token)
	if err != nil {
		t.Fatal(err)
	}
	var vr struct {
		Valid  bool         `json:"valid"`
		Claims badge.Claims `json:"claims"`
	}
	_ = json.NewDecoder(resp2.Body).Decode(&vr)
	_ = resp2.Body.Close()
	if !vr.Valid || vr.Claims.OrgName != "Acme Recruiting" {
		t.Fatalf("verify wrong: %+v", vr)
	}

	// SVG endpoint serves SVG.
	resp3, err := http.Get(srv.URL + "/badge/" + ir.Token + "/svg")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp3.Body.Close() }()
	if ct := resp3.Header.Get("Content-Type"); ct != "image/svg+xml" {
		t.Errorf("svg content-type %q", ct)
	}
}

func TestBadge_InvalidRisk(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	r := chi.NewRouter()
	web.NewBadgeIssuer(r, "https://j/", priv, pub, time.Hour)
	srv := httptest.NewServer(r)
	defer srv.Close()
	body, _ := json.Marshal(map[string]string{
		"org_id":          "o",
		"verification_id": "v",
		"risk":            "purple",
	})
	resp, err := http.Post(srv.URL+"/api/v1/recruiter/badges", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("got %d", resp.StatusCode)
	}
}
