package nominatim_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/nominatim"
)

func TestSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"place_id":1,"osm_type":"node","osm_id":42,"lat":"25.2","lon":"55.27","display_name":"Burj Khalifa, Dubai","address":{"city":"Dubai","country_code":"ae"}}]`))
	}))
	defer srv.Close()

	c := nominatim.New(srv.URL)
	places, err := c.Search(context.Background(), "burj khalifa", "ae", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 || places[0].Address.City != "Dubai" {
		t.Fatalf("unexpected: %+v", places)
	}
	lat, lon, err := places[0].LatLon()
	if err != nil || lat == 0 || lon == 0 {
		t.Errorf("lat,lon = %v,%v err=%v", lat, lon, err)
	}
}

func TestReverse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"place_id":2,"display_name":"Marina, Dubai","address":{"city":"Dubai","country_code":"ae"}}`))
	}))
	defer srv.Close()
	c := nominatim.New(srv.URL)
	p, err := c.Reverse(context.Background(), 25.0769, 55.1397)
	if err != nil || p == nil {
		t.Fatalf("p=%v err=%v", p, err)
	}
	if p.Address.CountryCode != "ae" {
		t.Errorf("country=%q", p.Address.CountryCode)
	}
}

func TestRateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	c := nominatim.New(srv.URL)
	_, err := c.Search(context.Background(), "x", "", 1)
	if !errors.Is(err, nominatim.ErrRateLimited) {
		t.Fatalf("want ErrRateLimited, got %v", err)
	}
}
