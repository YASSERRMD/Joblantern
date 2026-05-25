package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/nominatim"
	"github.com/yasserrmd/joblantern/internal/overpass"
)

// Structured error codes surfaced to MCP consumers.
const (
	ErrAddressNotFound = "ADDRESS_NOT_FOUND"
	ErrRateLimited     = "RATE_LIMITED"
	ErrUpstreamTimeout = "UPSTREAM_TIMEOUT"
	ErrInvalidArgs     = "INVALID_ARGS"
)

// ---- Tool: verify_address_exists ----

type verifyArgs struct {
	Address     string `json:"address" jsonschema:"freeform postal address"`
	CountryCode string `json:"country_code,omitempty" jsonschema:"ISO 3166-1 alpha-2 hint"`
}

type verifyResult struct {
	Found       bool    `json:"found"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	OSMType     string  `json:"osm_type,omitempty"`
	OSMID       int64   `json:"osm_id,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	Code        string  `json:"code,omitempty"`
}

// ---- Tool: reverse_geocode ----

type reverseArgs struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type reverseResult struct {
	Found       bool   `json:"found"`
	DisplayName string `json:"display_name,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Code        string `json:"code,omitempty"`
}

// ---- Tool: classify_building_type ----

type classifyArgs struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m,omitempty"`
}

type classifyResult struct {
	overpass.BuildingClassification
	Code string `json:"code,omitempty"`
}

// ---- Tool: address_cluster_check ----
// Wired to the local scam DB by the agent later. The MCP server returns
// a stub of zero when no DB backend is configured.

type clusterArgs struct {
	Address string `json:"address"`
}

type clusterResult struct {
	Count int    `json:"count"`
	Code  string `json:"code,omitempty"`
}

// ---- Tool: _meta_attribution ----
// MCP tool names cannot contain '/', so we expose this as _meta_attribution.

type attribArgs struct{}
type attribResult struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func addTools(s *mcp.Server, d deps) {
	mcp.AddTool(s,
		&mcp.Tool{Name: "verify_address_exists", Description: "Forward geocode a postal address via Nominatim."},
		func(ctx context.Context, _ *mcp.CallToolRequest, args verifyArgs) (*mcp.CallToolResult, verifyResult, error) {
			if args.Address == "" {
				return errCT(ErrInvalidArgs, "address is required"), verifyResult{Code: ErrInvalidArgs}, nil
			}
			places, err := d.nom.Search(ctx, args.Address, args.CountryCode, 1)
			if err != nil {
				if errors.Is(err, nominatim.ErrRateLimited) {
					return errCT(ErrRateLimited, err.Error()), verifyResult{Code: ErrRateLimited}, nil
				}
				return errCT(ErrUpstreamTimeout, err.Error()), verifyResult{Code: ErrUpstreamTimeout}, nil
			}
			if len(places) == 0 {
				return okCT("no match"), verifyResult{Found: false, Code: ErrAddressNotFound}, nil
			}
			p := places[0]
			lat, lon, _ := p.LatLon()
			return okCT(p.DisplayName), verifyResult{
				Found: true, Lat: lat, Lon: lon,
				OSMType: p.OSMType, OSMID: p.OSMID, DisplayName: p.DisplayName,
			}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "reverse_geocode", Description: "Reverse geocode a coordinate via Nominatim."},
		func(ctx context.Context, _ *mcp.CallToolRequest, args reverseArgs) (*mcp.CallToolResult, reverseResult, error) {
			p, err := d.nom.Reverse(ctx, args.Lat, args.Lon)
			if err != nil {
				if errors.Is(err, nominatim.ErrRateLimited) {
					return errCT(ErrRateLimited, err.Error()), reverseResult{Code: ErrRateLimited}, nil
				}
				return errCT(ErrUpstreamTimeout, err.Error()), reverseResult{Code: ErrUpstreamTimeout}, nil
			}
			if p == nil {
				return okCT("no match"), reverseResult{Found: false, Code: ErrAddressNotFound}, nil
			}
			return okCT(p.DisplayName), reverseResult{
				Found: true, DisplayName: p.DisplayName,
				CountryCode: p.Address.CountryCode, City: p.Address.City, State: p.Address.State,
			}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "classify_building_type", Description: "Classify the dominant land use within radius_m around (lat, lon)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, args classifyArgs) (*mcp.CallToolResult, classifyResult, error) {
			radius := args.RadiusM
			if radius <= 0 || radius > 500 {
				radius = 50
			}
			r, err := d.over.NearbyFeatures(ctx, args.Lat, args.Lon, radius)
			if err != nil {
				if errors.Is(err, overpass.ErrRateLimited) {
					return errCT(ErrRateLimited, err.Error()), classifyResult{Code: ErrRateLimited}, nil
				}
				return errCT(ErrUpstreamTimeout, err.Error()), classifyResult{Code: ErrUpstreamTimeout}, nil
			}
			c := overpass.Classify(r)
			return okCT(c.PrimaryType), classifyResult{BuildingClassification: c}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "address_cluster_check", Description: "Count how many distinct companies in the scam DB share this address. Stub in v1."},
		func(_ context.Context, _ *mcp.CallToolRequest, _ clusterArgs) (*mcp.CallToolResult, clusterResult, error) {
			return okCT("0"), clusterResult{Count: 0}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "_meta_attribution", Description: "Required OpenStreetMap attribution string for any UI displaying these results."},
		func(_ context.Context, _ *mcp.CallToolRequest, _ attribArgs) (*mcp.CallToolResult, attribResult, error) {
			r := attribResult{
				Text: "© OpenStreetMap contributors",
				URL:  "https://www.openstreetmap.org/copyright",
			}
			return okCT(r.Text), r, nil
		})
}

func okCT(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}

func errCT(code, msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("[%s] %s", code, msg)}},
	}
}
