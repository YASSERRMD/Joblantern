// Package vies is a tiny client for the EU's VIES VAT validation REST
// endpoint (https://ec.europa.eu/taxation_customs/vies/rest-api/).
// We use it to confirm an EU-jurisdiction recruiter's claimed VAT id
// is real before the agent treats the recruiter as legitimate.
package vies

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBase = "https://ec.europa.eu/taxation_customs/vies/rest-api"

var (
	ErrInvalidCountry = errors.New("vies: invalid country code")
	ErrNotFound       = errors.New("vies: VAT id not found")
	ErrUpstream       = errors.New("vies: upstream error")
)

// Result is the validation outcome.
type Result struct {
	Valid     bool      `json:"valid"`
	Country   string    `json:"country"`
	VATNumber string    `json:"vat_number"`
	Name      string    `json:"name,omitempty"`
	Address   string    `json:"address,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// Client talks to VIES.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

func New() *Client {
	return &Client{
		BaseURL:    defaultBase,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		UserAgent:  "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
	}
}

// Validate checks (countryCode, vatNumber). countryCode is the
// two-letter ISO-3166-1 alpha-2 code (e.g. "DE", "FR", "IT").
func (c *Client) Validate(ctx context.Context, countryCode, vatNumber string) (*Result, error) {
	cc := strings.ToUpper(strings.TrimSpace(countryCode))
	if len(cc) != 2 {
		return nil, ErrInvalidCountry
	}
	vn := strings.TrimSpace(vatNumber)
	if vn == "" {
		return nil, fmt.Errorf("vies: empty vat number")
	}
	endpoint := fmt.Sprintf("%s/check-vat-number/%s/%s", c.BaseURL, cc, vn)
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
	case http.StatusNotFound:
		return nil, ErrNotFound
	default:
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("%w: %d %s", ErrUpstream, resp.StatusCode, string(raw))
	}

	var body struct {
		IsValid bool   `json:"isValid"`
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("vies decode: %w", err)
	}
	return &Result{
		Valid:     body.IsValid,
		Country:   cc,
		VATNumber: vn,
		Name:      body.Name,
		Address:   body.Address,
		CheckedAt: time.Now().UTC(),
	}, nil
}
