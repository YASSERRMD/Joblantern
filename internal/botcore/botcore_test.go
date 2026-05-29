package botcore_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/botcore"
)

func TestAPI_VerifyAndWait(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/verify":
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]string{"verification_id": "abc"})
		case strings.HasPrefix(r.URL.Path, "/api/v1/verifications/"):
			calls++
			rec := botcore.Record{ID: "abc", Status: "running"}
			if calls > 1 {
				rec.Status = "completed"
				rec.Verdict = &botcore.Verdict{
					OverallRisk: "red",
					Confidence:  0.85,
					Reasons:     []string{"upfront_fee"},
				}
			}
			_ = json.NewEncoder(w).Encode(rec)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c := botcore.New(srv.URL, "")
	id, err := c.Verify(context.Background(), botcore.Submission{ListingText: "x", Jurisdiction: "AE"})
	if err != nil || id != "abc" {
		t.Fatalf("verify id=%q err=%v", id, err)
	}
	rec, err := c.Wait(context.Background(), id, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status != "completed" || rec.Verdict.OverallRisk != "red" {
		t.Fatalf("got %+v", rec)
	}
}

func TestRateLimiter(t *testing.T) {
	rl := botcore.NewRateLimiter(2, time.Minute)
	first, second := rl.Allow("a"), rl.Allow("a")
	if !first || !second {
		t.Fatal("first two should pass")
	}
	if rl.Allow("a") {
		t.Fatal("third should be blocked")
	}
	if !rl.Allow("b") {
		t.Fatal("other key should be independent")
	}
}

func TestFormatVerdict(t *testing.T) {
	out := botcore.FormatVerdict(&botcore.Record{
		ID: "abc", Status: "completed",
		Verdict: &botcore.Verdict{OverallRisk: "red", Confidence: 0.9, Reasons: []string{"upfront_fee"}},
	}, "https://example/verifications/abc")
	if !strings.Contains(out, "RED") || !strings.Contains(out, "90%") {
		t.Errorf("missing risk/conf: %q", out)
	}
}
