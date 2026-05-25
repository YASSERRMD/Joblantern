// Package nominatim is a small client for a self-hosted Nominatim
// instance. It deliberately exposes only the endpoints mcp-address
// needs: /search and /reverse.
package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client talks to a Nominatim HTTP service.
type Client struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

// New returns a Client with sensible defaults. baseURL has no trailing
// slash (e.g. "http://localhost:8088").
func New(baseURL string) *Client {
	return &Client{
		BaseURL:   baseURL,
		UserAgent: "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Place is a subset of the Nominatim JSON response we care about.
type Place struct {
	PlaceID     int64   `json:"place_id"`
	OSMType     string  `json:"osm_type"`
	OSMID       int64   `json:"osm_id"`
	Lat         string  `json:"lat"`
	Lon         string  `json:"lon"`
	DisplayName string  `json:"display_name"`
	Class       string  `json:"class"`
	Type        string  `json:"type"`
	Importance  float64 `json:"importance"`
	Address     Address `json:"address"`
}

// Address mirrors Nominatim's address object (most fields optional).
type Address struct {
	HouseNumber string `json:"house_number,omitempty"`
	Road        string `json:"road,omitempty"`
	Suburb      string `json:"suburb,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Postcode    string `json:"postcode,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

// LatLon returns the place's coordinates as floats. Returns 0,0,error
// if either field is malformed.
func (p Place) LatLon() (lat, lon float64, err error) {
	lat, err = strconv.ParseFloat(p.Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lat: %w", err)
	}
	lon, err = strconv.ParseFloat(p.Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse lon: %w", err)
	}
	return lat, lon, nil
}

// Search performs a forward geocode. countryCode is optional ("ae", "in", …).
func (c *Client) Search(ctx context.Context, query, countryCode string, limit int) ([]Place, error) {
	if limit <= 0 {
		limit = 5
	}
	q := url.Values{}
	q.Set("q", query)
	q.Set("format", "jsonv2")
	q.Set("addressdetails", "1")
	q.Set("limit", strconv.Itoa(limit))
	if countryCode != "" {
		q.Set("countrycodes", countryCode)
	}
	return c.getPlaces(ctx, "/search?"+q.Encode())
}

// Reverse performs a reverse geocode.
func (c *Client) Reverse(ctx context.Context, lat, lon float64) (*Place, error) {
	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(lat, 'f', 7, 64))
	q.Set("lon", strconv.FormatFloat(lon, 'f', 7, 64))
	q.Set("format", "jsonv2")
	q.Set("addressdetails", "1")
	places, err := c.getPlaces(ctx, "/reverse?"+q.Encode())
	if err != nil {
		return nil, err
	}
	if len(places) == 0 {
		return nil, nil
	}
	return &places[0], nil
}

func (c *Client) getPlaces(ctx context.Context, path string) ([]Place, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
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

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("nominatim %s: %d %s", path, resp.StatusCode, string(body))
	}

	// /search returns []Place, /reverse returns a single Place object.
	// Normalise both to []Place.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var arr []Place
	if err := json.Unmarshal(data, &arr); err == nil {
		return arr, nil
	}
	var single Place
	if err := json.Unmarshal(data, &single); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if single.PlaceID == 0 && single.DisplayName == "" {
		return nil, nil
	}
	return []Place{single}, nil
}
