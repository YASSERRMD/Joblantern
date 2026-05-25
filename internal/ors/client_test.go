package ors_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/ors"
)

func TestRoute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"routes":[{"summary":{"distance":15234.7,"duration":1240.0}}]}`))
	}))
	defer srv.Close()
	c := ors.New("key")
	c.BaseURL = srv.URL
	r, err := c.Route(context.Background(), ors.ModeDriving, 25.0, 55.0, 25.3, 55.4)
	if err != nil {
		t.Fatal(err)
	}
	if r.DistanceM <= 0 || r.DurationS <= 0 {
		t.Fatalf("got %+v", r)
	}
}

func TestRoute_OutOfRegion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Could not find routable point"}`))
	}))
	defer srv.Close()
	c := ors.New("key")
	c.BaseURL = srv.URL
	_, err := c.Route(context.Background(), ors.ModeDriving, 0, 0, 0, 0)
	if !errors.Is(err, ors.ErrOutOfRegion) {
		t.Fatalf("want ErrOutOfRegion got %v", err)
	}
}
