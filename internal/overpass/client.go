// Package overpass is a small client for an Overpass API instance.
// We use it to inspect land-use / building / shop / office tags near a
// coordinate so the address MCP server can label an address as
// residential, commercial, or mixed.
package overpass

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client talks to an Overpass HTTP service.
type Client struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL:   baseURL,
		UserAgent: "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Element is the subset of Overpass JSON we use.
type Element struct {
	Type string            `json:"type"`
	ID   int64             `json:"id"`
	Tags map[string]string `json:"tags,omitempty"`
}

// Response is what /api/interpreter returns.
type Response struct {
	Elements []Element `json:"elements"`
}

// ErrRateLimited signals upstream throttling.
var ErrRateLimited = errors.New("overpass: rate limited")

// Query runs an Overpass QL query and returns the decoded response.
func (c *Client) Query(ctx context.Context, ql string) (*Response, error) {
	form := url.Values{}
	form.Set("data", ql)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL,
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusGatewayTimeout {
		return nil, ErrRateLimited
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("overpass: %d %s", resp.StatusCode, string(body))
	}

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &r, nil
}

// NearbyFeatures runs an Overpass query that fetches all building,
// landuse, office and shop tags within radius_m metres of (lat, lon)
// using the [out:json] header and a single around: filter.
func (c *Client) NearbyFeatures(ctx context.Context, lat, lon float64, radiusM int) (*Response, error) {
	ql := fmt.Sprintf(`
[out:json][timeout:25];
(
  nwr["building"](around:%d,%f,%f);
  nwr["landuse"](around:%d,%f,%f);
  nwr["office"](around:%d,%f,%f);
  nwr["shop"](around:%d,%f,%f);
);
out tags 200;
`,
		radiusM, lat, lon,
		radiusM, lat, lon,
		radiusM, lat, lon,
		radiusM, lat, lon)
	return c.Query(ctx, ql)
}
