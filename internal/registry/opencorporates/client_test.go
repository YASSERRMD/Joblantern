package opencorporates_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/registry"
	"github.com/yasserrmd/joblantern/internal/registry/opencorporates"
)

func TestLookupByName_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":{"companies":[{"company":{"name":"Acme Ltd","company_number":"12345","jurisdiction_code":"gb","current_status":"Active","incorporation_date":"2020-01-15","opencorporates_url":"https://example/acme"}}]}}`))
	}))
	defer srv.Close()

	p := opencorporates.New("")
	p.BaseURL = srv.URL
	ms, err := p.LookupByName(context.Background(), "acme", "gb", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) != 1 || ms[0].Name != "Acme Ltd" || ms[0].Jurisdiction != "gb" {
		t.Fatalf("got %+v", ms)
	}
}

func TestLookupByName_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":{"companies":[]}}`))
	}))
	defer srv.Close()
	p := opencorporates.New("")
	p.BaseURL = srv.URL
	_, err := p.LookupByName(context.Background(), "ghost", "", 5)
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("want ErrNotFound got %v", err)
	}
}

func TestGet_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":{"company":{"name":"X","company_number":"1","jurisdiction_code":"ae","registered_address":{"in_full":"Dubai"},"officers":[{"officer":{"name":"Jane","position":"Director","start_date":"2024-01-01"}}]}}}`))
	}))
	defer srv.Close()
	p := opencorporates.New("")
	p.BaseURL = srv.URL
	c, err := p.Get(context.Background(), "ae/1")
	if err != nil {
		t.Fatal(err)
	}
	if c.RegisteredAddress != "Dubai" || len(c.Officers) != 1 || c.Officers[0].Name != "Jane" {
		t.Fatalf("got %+v", c)
	}
}

func TestRateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	p := opencorporates.New("")
	p.BaseURL = srv.URL
	_, err := p.LookupByName(context.Background(), "x", "", 1)
	if !errors.Is(err, registry.ErrRateLimited) {
		t.Fatalf("want ErrRateLimited got %v", err)
	}
}
