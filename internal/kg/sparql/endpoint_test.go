package sparql

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRejectsEmptyQuery(t *testing.T) {
	srv := httptest.NewServer(Endpoint{Execute: func(_ string) ([]map[string]string, error) { return nil, nil }})
	defer srv.Close()
	resp, _ := http.Get(srv.URL)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestStableSelectShape(t *testing.T) {
	srv := httptest.NewServer(Endpoint{Execute: func(_ string) ([]map[string]string, error) {
		return []map[string]string{{"company": "ex:c1", "industry": "transport"}}, nil
	}})
	defer srv.Close()
	resp, _ := http.Get(srv.URL + "?query=SELECT+%3Fcompany+WHERE+%7B+%7D")
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "bindings") {
		t.Fatalf("expected sparql-results json, got %s", string(body))
	}
}
