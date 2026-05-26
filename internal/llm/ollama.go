// Package llm hosts provider adapters for local and remote LLM
// backends. Phase 13 wired the agent's ScoreFunc; this package fills
// in the natural-language narrative around the deterministic verdict.
//
// v0.1 ships an Ollama adapter (local-first). Future PRs add llama.cpp
// HTTP, OpenAI-compatible endpoints, and a Bifrost-style multi-provider
// router. All adapters satisfy the same Provider interface so the
// agent can pick at runtime via config/llm.yaml.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Provider is the narrow surface every LLM backend implements.
type Provider interface {
	// Name is a stable identifier ("ollama", "anthropic", "openai").
	Name() string
	// Generate returns a single chat completion. Streaming is left
	// to a future revision once the UI surfaces streamed narrative.
	Generate(ctx context.Context, req Request) (Response, error)
}

// Request is the provider-neutral generation request.
type Request struct {
	Model       string
	System      string
	Prompt      string
	Temperature float64
	MaxTokens   int
}

// Response is the provider-neutral generation result.
type Response struct {
	Text   string
	Tokens int // total tokens (input + output) where the backend reports it; 0 otherwise.
}

// OllamaProvider talks to a local Ollama server (default
// http://localhost:11434).
type OllamaProvider struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewOllama returns a Provider pointing at baseURL (empty = default).
func NewOllama(baseURL string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OllamaProvider) Name() string { return "ollama" }

type ollamaReq struct {
	Model   string         `json:"model"`
	System  string         `json:"system,omitempty"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options,omitempty"`
}

type ollamaResp struct {
	Response   string `json:"response"`
	PromptEval int    `json:"prompt_eval_count,omitempty"`
	Eval       int    `json:"eval_count,omitempty"`
}

// Generate hits the /api/generate endpoint with stream=false.
func (p *OllamaProvider) Generate(ctx context.Context, req Request) (Response, error) {
	if req.Model == "" {
		return Response{}, errors.New("ollama: model is required")
	}
	options := map[string]any{}
	if req.Temperature > 0 {
		options["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		options["num_predict"] = req.MaxTokens
	}
	body, _ := json.Marshal(ollamaReq{
		Model:   req.Model,
		System:  req.System,
		Prompt:  req.Prompt,
		Stream:  false,
		Options: options,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.BaseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return Response{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return Response{}, fmt.Errorf("ollama: %d %s", resp.StatusCode, string(raw))
	}
	var out ollamaResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Response{}, fmt.Errorf("ollama decode: %w", err)
	}
	return Response{
		Text:   out.Response,
		Tokens: out.PromptEval + out.Eval,
	}, nil
}
