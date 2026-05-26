// Package companieshouse implements registry.Provider against the
// UK Companies House public API.
//
// Auth: HTTP Basic with the API key as username, empty password.
// Docs: https://developer-specs.company-information.service.gov.uk/
package companieshouse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yasserrmd/joblantern/internal/registry"
)

const defaultBase = "https://api.company-information.service.gov.uk"

// Provider talks to Companies House.
type Provider struct {
	BaseURL    string
	APIKey     string
	UserAgent  string
	HTTPClient *http.Client
}

func New(apiKey string) *Provider {
	return &Provider{
		BaseURL:    defaultBase,
		APIKey:     apiKey,
		UserAgent:  "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *Provider) Name() string { return "companieshouse-uk" }

type chSearchHit struct {
	Title          string `json:"title"`
	CompanyNumber  string `json:"company_number"`
	CompanyStatus  string `json:"company_status"`
	DateOfCreation string `json:"date_of_creation"`
	Links          struct {
		Self string `json:"self"`
	} `json:"links"`
}

type chSearchEnvelope struct {
	Items []chSearchHit `json:"items"`
}

type chCompany struct {
	CompanyName       string `json:"company_name"`
	CompanyNumber     string `json:"company_number"`
	CompanyStatus     string `json:"company_status"`
	DateOfCreation    string `json:"date_of_creation"`
	RegisteredAddress struct {
		AddressLine1 string `json:"address_line_1"`
		Locality     string `json:"locality"`
		PostalCode   string `json:"postal_code"`
		Country      string `json:"country"`
	} `json:"registered_office_address"`
}

// LookupByName satisfies registry.Provider. `jurisdiction` is ignored
// because Companies House only knows UK companies.
func (p *Provider) LookupByName(ctx context.Context, name, _ string, limit int) ([]registry.Match, error) {
	if limit <= 0 {
		limit = 5
	}
	q := url.Values{}
	q.Set("q", name)
	q.Set("items_per_page", strconv.Itoa(limit))
	endpoint := p.BaseURL + "/search/companies?" + q.Encode()

	var env chSearchEnvelope
	if err := p.do(ctx, endpoint, &env); err != nil {
		return nil, err
	}
	out := make([]registry.Match, 0, len(env.Items))
	for _, it := range env.Items {
		out = append(out, registry.Match{
			ID:                "gb/" + it.CompanyNumber,
			Name:              it.Title,
			Jurisdiction:      "gb",
			Status:            it.CompanyStatus,
			IncorporationDate: parseDate(it.DateOfCreation),
			RegistryURL:       "https://find-and-update.company-information.service.gov.uk/company/" + it.CompanyNumber,
		})
	}
	if len(out) == 0 {
		return nil, registry.ErrNotFound
	}
	return out, nil
}

// Get satisfies registry.Provider. id is "gb/<company-number>".
func (p *Provider) Get(ctx context.Context, id string) (*registry.Company, error) {
	const prefix = "gb/"
	if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
		return nil, registry.ErrNotFound
	}
	num := id[len(prefix):]
	endpoint := p.BaseURL + "/company/" + url.PathEscape(num)
	var c chCompany
	if err := p.do(ctx, endpoint, &c); err != nil {
		return nil, err
	}
	return &registry.Company{
		Match: registry.Match{
			ID: id, Name: c.CompanyName, Jurisdiction: "gb",
			Status:            c.CompanyStatus,
			IncorporationDate: parseDate(c.DateOfCreation),
			RegistryURL:       "https://find-and-update.company-information.service.gov.uk/company/" + c.CompanyNumber,
		},
		RegisteredAddress: oneLine(c.RegisteredAddress.AddressLine1, c.RegisteredAddress.Locality, c.RegisteredAddress.PostalCode, c.RegisteredAddress.Country),
	}, nil
}

func (p *Provider) do(ctx context.Context, endpoint string, into any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", p.UserAgent)
	req.Header.Set("Accept", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(p.APIKey+":")))
	}
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return registry.ErrNotFound
	case http.StatusTooManyRequests:
		return registry.ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return registry.ErrTokenInvalid
	default:
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("companies-house: %d %s", resp.StatusCode, string(raw))
	}
	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
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

func oneLine(parts ...string) string {
	out := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		if out != "" {
			out += ", "
		}
		out += p
	}
	return out
}
