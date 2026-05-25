package mapillary_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/mapillary"
)

func TestImagesNearPoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/images" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"42","captured_at":1700000000000,"thumb_1024_url":"https://example/x.jpg"}]}`))
	}))
	defer srv.Close()
	c := mapillary.New("token")
	c.BaseURL = srv.URL
	imgs, err := c.ImagesNearPoint(context.Background(), 25.2, 55.27, 50, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(imgs) != 1 || imgs[0].ID != "42" {
		t.Fatalf("got %+v", imgs)
	}
	if imgs[0].CapturedTime().IsZero() {
		t.Error("expected non-zero CapturedTime")
	}
}

func TestTokenInvalid(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	c := mapillary.New("bad")
	c.BaseURL = srv.URL
	_, err := c.ImagesNearPoint(context.Background(), 0, 0, 50, 1)
	if !errors.Is(err, mapillary.ErrTokenInvalid) {
		t.Fatalf("want ErrTokenInvalid got %v", err)
	}
}
