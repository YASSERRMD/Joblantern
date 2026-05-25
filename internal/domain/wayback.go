package domain

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

// ArchiveSummary is the derived Wayback view used by mcp-domain.
type ArchiveSummary struct {
	Domain            string    `json:"domain"`
	EarliestSnapshot  time.Time `json:"earliest_snapshot,omitempty"`
	LatestSnapshot    time.Time `json:"latest_snapshot,omitempty"`
	SnapshotCount     int       `json:"snapshot_count"`
}

// WaybackClient queries the Internet Archive Wayback CDX API.
type WaybackClient struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
}

func NewWaybackClient() *WaybackClient {
	return &WaybackClient{
		BaseURL:   "https://web.archive.org",
		UserAgent: "Joblantern/0.x (+https://github.com/yasserrmd/joblantern)",
		HTTPClient: &http.Client{Timeout: 20 * time.Second},
	}
}

// Summary queries the CDX API for the first/last/count snapshots of
// the domain.
func (w *WaybackClient) Summary(ctx context.Context, dom string) (*ArchiveSummary, error) {
	q := url.Values{}
	q.Set("url", dom)
	q.Set("output", "json")
	q.Set("fl", "timestamp")
	q.Set("filter", "statuscode:200")
	q.Set("collapse", "timestamp:8")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.BaseURL+"/cdx/search/cdx?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", w.UserAgent)
	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("wayback: %d %s", resp.StatusCode, string(body))
	}

	// CDX with output=json returns [["timestamp"], ["2010..."], ...]
	var rows [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	out := &ArchiveSummary{Domain: dom}
	for i, r := range rows {
		if i == 0 || len(r) == 0 {
			continue
		}
		ts := r[0]
		if _, err := strconv.ParseInt(ts, 10, 64); err != nil {
			continue
		}
		t, err := time.Parse("20060102150405", ts)
		if err != nil {
			continue
		}
		if out.EarliestSnapshot.IsZero() || t.Before(out.EarliestSnapshot) {
			out.EarliestSnapshot = t
		}
		if t.After(out.LatestSnapshot) {
			out.LatestSnapshot = t
		}
		out.SnapshotCount++
	}
	return out, nil
}
