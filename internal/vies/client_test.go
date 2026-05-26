package vies_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/vies"
)

func TestValidate_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"isValid":true,"name":"Acme GmbH","address":"Berlin"}`))
	}))
	defer srv.Close()
	c := vies.New()
	c.BaseURL = srv.URL
	r, err := c.Validate(context.Background(), "DE", "123456789")
	if err != nil || !r.Valid || r.Name != "Acme GmbH" {
		t.Fatalf("got %+v err=%v", r, err)
	}
}

func TestValidate_BadCountry(t *testing.T) {
	c := vies.New()
	_, err := c.Validate(context.Background(), "ZZZ", "1")
	if !errors.Is(err, vies.ErrInvalidCountry) {
		t.Fatalf("want ErrInvalidCountry got %v", err)
	}
}
