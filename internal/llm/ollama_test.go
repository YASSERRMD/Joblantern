package llm_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yasserrmd/joblantern/internal/llm"
)

func TestOllama_Generate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"response":          "Looks suspicious because of the upfront fee.",
			"prompt_eval_count": 20,
			"eval_count":        15,
		})
	}))
	defer srv.Close()

	p := llm.NewOllama(srv.URL)
	out, err := p.Generate(context.Background(), llm.Request{
		Model: "qwen2.5:3b", Prompt: "summarise verdict",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Text == "" || out.Tokens != 35 {
		t.Fatalf("got %+v", out)
	}
}

func TestOllama_MissingModel(t *testing.T) {
	p := llm.NewOllama("http://localhost:11434")
	_, err := p.Generate(context.Background(), llm.Request{Prompt: "x"})
	if err == nil {
		t.Fatal("expected error for missing model")
	}
}
