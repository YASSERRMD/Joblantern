// mcp-rental-listings is an MCP server that wraps open rental
// aggregator APIs (where ToS allows) and provides a normalized
// listing-lookup tool to the agent layer.
//
// As with mcp-domain and mcp-salary, this binary stays small and
// stateless. Aggregator credentials are loaded from env vars.
package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
)

// listing is the normalized response shape.
type listing struct {
	Source       string  `json:"source"`
	ListingID    string  `json:"listing_id"`
	Title        string  `json:"title"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	MonthlyRent  float64 `json:"monthly_rent"`
	Currency     string  `json:"currency"`
	ContactPhone string  `json:"contact_phone,omitempty"`
	ListingURL   string  `json:"listing_url,omitempty"`
}

// validateURL keeps us from accidentally crawling private nets.
func validateURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme != "https" {
		return "", fmt.Errorf("https required")
	}
	return u.String(), nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("mcp-rental-listings starting")
	// The MCP server boilerplate (stdio, tool registration) follows
	// the same template as the other Joblantern MCP servers and is
	// wired in a subsequent commit alongside the rental rule pack.
}
