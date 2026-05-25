package domain_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/domain"
)

func TestCrtSHSummary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[
			{"not_before":"2020-01-01T00:00:00","not_after":"2021-01-01T00:00:00","issuer_name":"Let's Encrypt"},
			{"not_before":"2024-05-01T00:00:00","not_after":"2024-08-01T00:00:00","issuer_name":"Let's Encrypt"}
		]`))
	}))
	defer srv.Close()
	c := domain.NewCrtSHClient()
	c.BaseURL = srv.URL
	s, err := c.Summary(context.Background(), "example.com")
	if err != nil {
		t.Fatal(err)
	}
	if s.CertCount != 2 || s.FirstCertAt.Year() != 2020 || len(s.UniqueIssuers) != 1 {
		t.Fatalf("%+v", s)
	}
}

func TestWaybackSummary(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[["timestamp"],["20100101000000"],["20240525000000"]]`))
	}))
	defer srv.Close()
	wc := domain.NewWaybackClient()
	wc.BaseURL = srv.URL
	a, err := wc.Summary(context.Background(), "example.com")
	if err != nil {
		t.Fatal(err)
	}
	if a.SnapshotCount != 2 || a.EarliestSnapshot.Year() != 2010 {
		t.Fatalf("%+v", a)
	}
}

// fakeWHOIS implements WHOISLookup with deterministic values.
type fakeWHOIS struct {
	w   *domain.WHOIS
	err error
}

func (f *fakeWHOIS) Lookup(_ context.Context, _ string) (*domain.WHOIS, error) {
	return f.w, f.err
}

func TestProfile_Composition(t *testing.T) {
	crt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"not_before":"2024-01-01T00:00:00","not_after":"2025-01-01T00:00:00","issuer_name":"X"}]`))
	}))
	defer crt.Close()
	wb := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[["timestamp"],["20231201000000"]]`))
	}))
	defer wb.Close()

	c := &domain.Composer{
		WHOIS:   &fakeWHOIS{w: &domain.WHOIS{Domain: "x.com", CreatedAt: time.Now().AddDate(-1, 0, 0)}},
		CrtSH:   &domain.CrtSHClient{BaseURL: crt.URL, HTTPClient: &http.Client{Timeout: 5 * time.Second}},
		Wayback: &domain.WaybackClient{BaseURL: wb.URL, HTTPClient: &http.Client{Timeout: 5 * time.Second}},
	}
	p, err := c.FullProfile(context.Background(), "x.com")
	if err != nil {
		t.Fatal(err)
	}
	if p.AgeDays < 300 {
		t.Errorf("age too short: %d", p.AgeDays)
	}
	if p.FreshnessScore == 1.0 || p.FreshnessScore == 0.0 {
		t.Errorf("expected mid freshness, got %v", p.FreshnessScore)
	}
}
