package overpass_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/overpass"
)

func TestNearbyFeatures(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "want post", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"elements":[
		  {"type":"way","id":1,"tags":{"building":"yes","office":"company"}},
		  {"type":"way","id":2,"tags":{"landuse":"residential"}}
		]}`))
	}))
	defer srv.Close()
	c := overpass.New(srv.URL)
	r, err := c.NearbyFeatures(context.Background(), 25.2, 55.27, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Elements) != 2 {
		t.Fatalf("got %d elements", len(r.Elements))
	}
}
