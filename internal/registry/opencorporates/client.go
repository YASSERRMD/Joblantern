// Package opencorporates implements registry.Provider using the
// OpenCorporates v0.4.8 REST API.
package opencorporates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yasserrmd/joblantern/internal/registry"
)

const defaultBase = "https://api.opencorporates.com/v0.4.8"

// Provider talks to OpenCorporates.
type Provider struct {
	BaseURL    string
	Token      string
	UserAgent  string
	HTTPClient *http.Client
}

// New returns a Provider with sensible defaults. token may be empty for
// the rate-limited anonymous tier.
func New(token string) *Provider {
	return &Provider{
		BaseURL:   defaultBase,
		Token:     token,
		UserAgent: "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *Provider) Name() string { return "opencorporates" }

// --- Wire types (subset) ---

type companyEnvelope struct {
	Results struct {
		Companies []companyHit `json:"companies"`
		Company   *companyRow  `json:"company,omitempty"`
	} `json:"results"`
}

type companyHit struct {
	Company companyRow `json:"company"`
}

type companyRow struct {
	Name              string `json:"name"`
	CompanyNumber     string `json:"company_number"`
	JurisdictionCode  string `json:"jurisdiction_code"`
	IncorporationDate string `json:"incorporation_date"`
	CurrentStatus     string `json:"current_status"`
	OpencorporatesURL string `json:"opencorporates_url"`
	RegisteredAddress struct {
		InFull string `json:"in_full"`
	} `json:"registered_address"`
	Officers []officerRow `json:"officers,omitempty"`
}

type officerRow struct {
	Officer struct {
		Name      string `json:"name"`
		Position  string `json:"position"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"officer"`
}

// LookupByName satisfies registry.Provider.
func (p *Provider) LookupByName(ctx context.Context, name, jurisdiction string, limit int) ([]registry.Match, error) {
	if limit <= 0 {
		limit = 5
	}
	q := url.Values{}
	q.Set("q", name)
	if jurisdiction != "" {
		q.Set("jurisdiction_code", jurisdiction)
	}
	q.Set("format", "json")
	q.Set("per_page", strconv.Itoa(limit))
	if p.Token != "" {
		q.Set("api_token", p.Token)
	}
	endpoint := p.BaseURL + "/companies/search?" + q.Encode()

	env, err := p.do(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	out := make([]registry.Match, 0, len(env.Results.Companies))
	for _, h := range env.Results.Companies {
		out = append(out, toMatch(h.Company))
	}
	if len(out) == 0 {
		return nil, registry.ErrNotFound
	}
	return out, nil
}

// Get satisfies registry.Provider. id is "<jurisdiction_code>/<company_number>".
func (p *Provider) Get(ctx context.Context, id string) (*registry.Company, error) {
	q := url.Values{}
	q.Set("format", "json")
	if p.Token != "" {
		q.Set("api_token", p.Token)
	}
	endpoint := p.BaseURL + "/companies/" + id + "?" + q.Encode()
	env, err := p.do(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if env.Results.Company == nil {
		return nil, registry.ErrNotFound
	}
	c := *env.Results.Company
	out := registry.Company{
		Match:             toMatch(c),
		RegisteredAddress: c.RegisteredAddress.InFull,
		Officers:          make([]registry.Officer, 0, len(c.Officers)),
	}
	for _, o := range c.Officers {
		out.Officers = append(out.Officers, registry.Officer{
			Name:      o.Officer.Name,
			Position:  o.Officer.Position,
			StartDate: parseDate(o.Officer.StartDate),
			EndDate:   parseDate(o.Officer.EndDate),
		})
	}
	return &out, nil
}

func toMatch(c companyRow) registry.Match {
	return registry.Match{
		ID:                c.JurisdictionCode + "/" + c.CompanyNumber,
		Name:              c.Name,
		Jurisdiction:      c.JurisdictionCode,
		Status:            c.CurrentStatus,
		IncorporationDate: parseDate(c.IncorporationDate),
		RegistryURL:       c.OpencorporatesURL,
	}
}

func parseDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t
	}
	return time.Time{}
}

func (p *Provider) do(ctx context.Context, endpoint string) (*companyEnvelope, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.UserAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, registry.ErrNotFound
	case http.StatusTooManyRequests:
		return nil, registry.ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, registry.ErrTokenInvalid
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("opencorporates: %d %s", resp.StatusCode, string(body))
	}
	var env companyEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &env, nil
}
