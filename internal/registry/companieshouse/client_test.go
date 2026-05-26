package companieshouse_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/registry/companieshouse"
)

func TestLookupAndGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/search/companies":
			_, _ = w.Write([]byte(`{"items":[{"title":"Acme Ltd","company_number":"12345","company_status":"active","date_of_creation":"2020-01-15"}]}`))
		case "/company/12345":
			_, _ = w.Write([]byte(`{"company_name":"Acme Ltd","company_number":"12345","company_status":"active","date_of_creation":"2020-01-15","registered_office_address":{"address_line_1":"1 High St","locality":"London","postal_code":"SW1A 1AA","country":"United Kingdom"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	p := companieshouse.New("k")
	p.BaseURL = srv.URL
	ms, err := p.LookupByName(context.Background(), "acme", "", 5)
	if err != nil || len(ms) != 1 || ms[0].ID != "gb/12345" {
		t.Fatalf("ms=%+v err=%v", ms, err)
	}
	c, err := p.Get(context.Background(), "gb/12345")
	if err != nil || c.RegisteredAddress == "" {
		t.Fatalf("c=%+v err=%v", c, err)
	}
}
