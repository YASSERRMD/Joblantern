// Package botcore is shared across messaging-channel adapters (Telegram,
// future WhatsApp, IVR). It owns the HTTP client that talks to a
// running Joblantern instance and the conversation state machine.
//
// Adapters supply transport-specific I/O (read a message, send a
// reply, fetch a file). The state machine is transport-agnostic.
package botcore

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

// APIClient talks to a Joblantern HTTP API.
type APIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// New returns an APIClient with sensible defaults.
func New(baseURL, apiKey string) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Submission mirrors agent.Submission. We re-declare it here so that
// the bot binary does not depend on internal/agent (those internals
// can churn). Only the JSON shape matters across the API boundary.
type Submission struct {
	ListingURL     string  `json:"listing_url,omitempty"`
	ListingText    string  `json:"listing_text,omitempty"`
	CompanyName    string  `json:"company_name,omitempty"`
	ClaimedAddress string  `json:"claimed_address,omitempty"`
	RecruiterEmail string  `json:"recruiter_email,omitempty"`
	RecruiterPhone string  `json:"recruiter_phone,omitempty"`
	Role           string  `json:"role,omitempty"`
	Jurisdiction   string  `json:"jurisdiction,omitempty"`
	Domain         string  `json:"domain,omitempty"`
	ClaimedSalary  float64 `json:"claimed_salary,omitempty"`
	SalaryCurrency string  `json:"salary_currency,omitempty"`
}

// Verdict is the subset of the server's verdict the bot renders.
type Verdict struct {
	OverallRisk string   `json:"overall_risk"`
	Confidence  float64  `json:"confidence"`
	Reasons     []string `json:"reasons"`
}

// Record matches the server's stored verification record.
type Record struct {
	ID      string   `json:"id"`
	Status  string   `json:"status"`
	Verdict *Verdict `json:"verdict,omitempty"`
}

// Verify POSTs the submission and returns the verification id.
func (c *APIClient) Verify(ctx context.Context, sub Submission) (string, error) {
	body, _ := json.Marshal(sub)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/v1/verify", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("X-Joblantern-API-Key", c.APIKey)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusAccepted {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("verify: %d %s", resp.StatusCode, string(raw))
	}
	var out struct {
		ID string `json:"verification_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.ID == "" {
		return "", errors.New("verify: empty id")
	}
	return out.ID, nil
}

// Get fetches one verification record.
func (c *APIClient) Get(ctx context.Context, id string) (*Record, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/v1/verifications/"+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("get: %d %s", resp.StatusCode, string(raw))
	}
	var rec Record
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

// Wait polls until the verification completes or times out.
func (c *APIClient) Wait(ctx context.Context, id string, max time.Duration) (*Record, error) {
	deadline := time.Now().Add(max)
	for time.Now().Before(deadline) {
		rec, err := c.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if rec.Status == "completed" || rec.Status == "failed" {
			return rec, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1500 * time.Millisecond):
		}
	}
	return nil, errors.New("verify: timeout")
}
