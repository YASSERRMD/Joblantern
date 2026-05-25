// Package mapillary is a small client for the Mapillary Graph API,
// scoped to the queries mcp-streetview needs.
//
// Auth: a client token from https://www.mapillary.com/dashboard/developers
// is required and is sent as a query parameter `access_token=...`.
package mapillary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultBase = "https://graph.mapillary.com"

// ErrRateLimited signals HTTP 429.
var ErrRateLimited = errors.New("mapillary: rate limited")

// ErrTokenInvalid signals HTTP 401/403.
var ErrTokenInvalid = errors.New("mapillary: token invalid")

// Client talks to the Mapillary Graph API.
type Client struct {
	BaseURL    string
	Token      string
	UserAgent  string
	HTTPClient *http.Client
}

func New(token string) *Client {
	return &Client{
		BaseURL:   defaultBase,
		Token:     token,
		UserAgent: "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Image is the subset of Mapillary's image entity we care about.
type Image struct {
	ID           string          `json:"id"`
	CapturedAt   int64           `json:"captured_at,omitempty"` // ms since epoch
	IsPano       bool            `json:"is_pano,omitempty"`
	SequenceID   string          `json:"sequence,omitempty"`
	ThumbURL     string          `json:"thumb_1024_url,omitempty"`
	OriginalURL  string          `json:"thumb_original_url,omitempty"`
	GeometryJSON json.RawMessage `json:"geometry,omitempty"`
}

// ImagesResponse wraps the Graph API's standard {data:[...]} envelope.
type ImagesResponse struct {
	Data []Image `json:"data"`
}

// ImagesNearPoint queries `images` filtered by `closeto` and `radius`.
// radiusM must be <= 50 m and limit <= 100 per Mapillary's docs.
func (c *Client) ImagesNearPoint(ctx context.Context, lat, lon float64, radiusM, limit int) ([]Image, error) {
	if radiusM <= 0 || radiusM > 50 {
		radiusM = 50
	}
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	q := url.Values{}
	q.Set("fields", "id,captured_at,is_pano,sequence,thumb_1024_url,thumb_original_url,geometry")
	q.Set("closeto", fmt.Sprintf("%f,%f", lon, lat)) // Graph expects lon,lat
	q.Set("radius", strconv.Itoa(radiusM))
	q.Set("limit", strconv.Itoa(limit))
	q.Set("access_token", c.Token)

	endpoint := c.BaseURL + "/images?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return nil, ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, ErrTokenInvalid
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("mapillary: %d %s", resp.StatusCode, string(body))
	}

	var ir ImagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return ir.Data, nil
}

// CapturedAt returns the image's capture time, or zero time if unset.
func (i Image) CapturedTime() time.Time {
	if i.CapturedAt == 0 {
		return time.Time{}
	}
	return time.UnixMilli(i.CapturedAt)
}
