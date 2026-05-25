package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// CertSummary is the derived SSL-history view used by mcp-domain.
type CertSummary struct {
	Domain        string    `json:"domain"`
	CertCount     int       `json:"cert_count"`
	FirstCertAt   time.Time `json:"first_cert_at,omitempty"`
	LastCertAt    time.Time `json:"last_cert_at,omitempty"`
	UniqueIssuers []string  `json:"unique_issuers"`
}

// CrtSHClient queries https://crt.sh.
type CrtSHClient struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

func NewCrtSHClient() *CrtSHClient {
	return &CrtSHClient{
		BaseURL:    "https://crt.sh",
		UserAgent:  "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
	}
}

type crtRow struct {
	NotBefore  string `json:"not_before"`
	NotAfter   string `json:"not_after"`
	IssuerName string `json:"issuer_name"`
}

// Summary returns a CertSummary for the given domain.
func (c *CrtSHClient) Summary(ctx context.Context, dom string) (*CertSummary, error) {
	q := url.Values{}
	q.Set("q", dom)
	q.Set("output", "json")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/?"+q.Encode(), nil)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("crt.sh: %d %s", resp.StatusCode, string(body))
	}
	var rows []crtRow
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	out := &CertSummary{Domain: dom, CertCount: len(rows)}
	seenIssuer := map[string]struct{}{}
	const fmtTS = "2006-01-02T15:04:05"
	for _, r := range rows {
		if t, err := time.Parse(fmtTS, r.NotBefore); err == nil {
			if out.FirstCertAt.IsZero() || t.Before(out.FirstCertAt) {
				out.FirstCertAt = t
			}
			if t.After(out.LastCertAt) {
				out.LastCertAt = t
			}
		}
		if _, dup := seenIssuer[r.IssuerName]; !dup && r.IssuerName != "" {
			seenIssuer[r.IssuerName] = struct{}{}
			out.UniqueIssuers = append(out.UniqueIssuers, r.IssuerName)
		}
	}
	return out, nil
}
