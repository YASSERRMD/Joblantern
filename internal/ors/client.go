// Package ors is a small client for the OpenRouteService API, scoped
// to the endpoints mcp-routing needs.
package ors

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

const defaultBase = "https://api.openrouteservice.org"

var (
	ErrRateLimited = errors.New("ors: rate limited")
	ErrOutOfRegion = errors.New("ors: out of region")
)

// Client talks to ORS.
type Client struct {
	BaseURL    string
	APIKey     string
	UserAgent  string
	HTTPClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		BaseURL:    defaultBase,
		APIKey:     apiKey,
		UserAgent:  "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
	}
}

// Mode is the ORS profile.
type Mode string

const (
	ModeDriving Mode = "driving-car"
	ModeWalking Mode = "foot-walking"
	ModeCycling Mode = "cycling-regular"
)

// RouteResult is a flattened subset of ORS directions output.
type RouteResult struct {
	DurationS float64 `json:"duration_s"`
	DistanceM float64 `json:"distance_m"`
}

// Route requests directions between two points.
func (c *Client) Route(ctx context.Context, mode Mode, fromLat, fromLon, toLat, toLon float64) (*RouteResult, error) {
	body := map[string]any{
		"coordinates": [][]float64{
			{fromLon, fromLat},
			{toLon, toLat},
		},
	}
	b, _ := json.Marshal(body)
	endpoint := fmt.Sprintf("%s/v2/directions/%s", c.BaseURL, mode)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode == http.StatusUnprocessableEntity || resp.StatusCode == http.StatusBadRequest {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if bytes.Contains(raw, []byte("Could not find routable point")) ||
			bytes.Contains(raw, []byte("not in any included regions")) {
			return nil, ErrOutOfRegion
		}
		return nil, fmt.Errorf("ors: %d %s", resp.StatusCode, string(raw))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("ors: %d %s", resp.StatusCode, string(raw))
	}
	var out struct {
		Routes []struct {
			Summary struct {
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
			} `json:"summary"`
		} `json:"routes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if len(out.Routes) == 0 {
		return nil, fmt.Errorf("ors: empty routes")
	}
	r := out.Routes[0]
	return &RouteResult{DurationS: r.Summary.Duration, DistanceM: r.Summary.Distance}, nil
}
