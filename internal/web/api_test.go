package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yasserrmd/joblantern/internal/agent"
	"github.com/yasserrmd/joblantern/internal/web"
)

type alwaysRed struct{}

func (alwaysRed) Name() string { return "always" }
func (alwaysRed) Run(_ context.Context, _ agent.Submission) []agent.Fact {
	return []agent.Fact{{
		Source: "test", ToolName: "t", FactType: "f",
		SupportsRisk: "red", Weight: 0.95,
	}}
}

func TestAPI_PostVerifyAndGet(t *testing.T) {
	r := chi.NewRouter()
	store := agent.NewMemoryStore()
	orch := agent.New(alwaysRed{})
	web.NewAPIHandler(r, store, orch)
	srv := httptest.NewServer(r)
	defer srv.Close()

	body, _ := json.Marshal(agent.Submission{ListingText: "x", Jurisdiction: "AE"})
	resp, err := http.Post(srv.URL+"/api/v1/verify", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("post status=%d body=%s", resp.StatusCode, raw)
	}
	var out map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&out)
	_ = resp.Body.Close()
	id := out["verification_id"]
	if id == "" {
		t.Fatal("no id returned")
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		r2, err := http.Get(srv.URL + "/api/v1/verifications/" + id)
		if err != nil {
			t.Fatal(err)
		}
		var rec agent.Record
		_ = json.NewDecoder(r2.Body).Decode(&rec)
		_ = r2.Body.Close()
		if rec.Status == "completed" {
			if rec.Verdict == nil {
				t.Fatal("completed but no verdict")
			}
			if rec.Verdict.OverallRisk != "yellow" && rec.Verdict.OverallRisk != "red" {
				t.Errorf("unexpected risk %q", rec.Verdict.OverallRisk)
			}
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("did not complete within deadline")
}
