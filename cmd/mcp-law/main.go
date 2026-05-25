// Package main is the joblantern.law MCP server.
//
// All rules are JSON-bundled. Future versions will hot-reload from
// data/recruitment-law.json and a database table (see migration 0009).
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerName = "joblantern.law"

const (
	ErrJurisdictionUnknown = "JURISDICTION_UNKNOWN"
	ErrCitationMissing     = "CITATION_MISSING"
	ErrNotImplemented      = "NOT_IMPLEMENTED"
)

//go:embed law.json
var lawJSON []byte

// Jurisdiction is one country's rule set.
type Jurisdiction struct {
	CountryCode         string `json:"country_code"`
	Name                string `json:"name"`
	RecruitmentFeeLegal bool   `json:"recruitment_fee_legal"`
	MaxFeePct           *int   `json:"max_fee_pct"`
	EmployerPaysOnly    bool   `json:"employer_pays_only"`
	LicensingRequired   bool   `json:"licensing_required"`
	RegulatorURL        string `json:"regulator_url"`
	CitationURL         string `json:"citation_url"`
}

type pack struct {
	Disclaimer    string         `json:"disclaimer"`
	Jurisdictions []Jurisdiction `json:"jurisdictions"`
}

func loadPack() (*pack, error) {
	var p pack
	if err := json.Unmarshal(lawJSON, &p); err != nil {
		return nil, err
	}
	for _, j := range p.Jurisdictions {
		if j.CountryCode == "" || j.CitationURL == "" {
			return nil, fmt.Errorf("jurisdiction missing required field: %+v", j)
		}
	}
	return &p, nil
}

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8089"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	p, err := loadPack()
	if err != nil {
		logger.Error("load law pack", "err", err)
		os.Exit(1)
	}

	s := newServer(p)
	switch *transport {
	case "stdio":
		if err := runStdio(s); err != nil {
			logger.Error("stdio", "err", err)
			os.Exit(1)
		}
	case "http":
		if err := runHTTP(s, *addr); err != nil {
			logger.Error("http", "err", err)
			os.Exit(1)
		}
	}
}

type feeArgs struct {
	Country         string  `json:"country"`
	ClaimedFee      float64 `json:"claimed_fee_amount"`
	ClaimedCurrency string  `json:"currency,omitempty"`
}
type feeResult struct {
	Legal        bool   `json:"legal"`
	Reason       string `json:"reason"`
	RegulatorURL string `json:"regulator_url,omitempty"`
	CitationURL  string `json:"citation_url,omitempty"`
	Code         string `json:"code,omitempty"`
}

type licenseArgs struct {
	Country       string `json:"country"`
	RecruiterName string `json:"recruiter_name,omitempty"`
	RecruiterID   string `json:"recruiter_id,omitempty"`
}
type licenseResult struct {
	RegulatorURL string `json:"regulator_url"`
	Note         string `json:"note"`
	Code         string `json:"code,omitempty"`
}

type visaArgs struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Role        string `json:"role,omitempty"`
}
type visaResult struct {
	CitationURL string `json:"citation_url,omitempty"`
	Note        string `json:"note"`
	Code        string `json:"code,omitempty"`
}

type disclaimerArgs struct{}
type disclaimerResult struct {
	Text string `json:"text"`
}

func newServer(p *pack) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "fee_legality_check", Description: "Is the claimed recruitment fee legal in the destination country?"},
		func(_ context.Context, _ *mcp.CallToolRequest, a feeArgs) (*mcp.CallToolResult, feeResult, error) {
			j := p.find(a.Country)
			if j == nil {
				return errCT(ErrJurisdictionUnknown, a.Country), feeResult{Code: ErrJurisdictionUnknown}, nil
			}
			r := feeResult{
				RegulatorURL: j.RegulatorURL, CitationURL: j.CitationURL,
			}
			if !j.RecruitmentFeeLegal {
				r.Legal = false
				r.Reason = fmt.Sprintf("Charging the worker a recruitment fee is illegal in %s; employer must pay.", j.Name)
			} else if j.MaxFeePct != nil {
				r.Legal = true
				r.Reason = fmt.Sprintf("Fee permitted up to %d%% of first-month salary in %s.", *j.MaxFeePct, j.Name)
			} else {
				r.Legal = true
				r.Reason = fmt.Sprintf("Fee permitted in %s; no statutory cap recorded — verify with the regulator.", j.Name)
			}
			return okCT(r.Reason), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "check_recruiter_license", Description: "v1 returns NOT_IMPLEMENTED + the regulator URL where a license can be verified."},
		func(_ context.Context, _ *mcp.CallToolRequest, a licenseArgs) (*mcp.CallToolResult, licenseResult, error) {
			j := p.find(a.Country)
			if j == nil {
				return errCT(ErrJurisdictionUnknown, a.Country), licenseResult{Code: ErrJurisdictionUnknown}, nil
			}
			r := licenseResult{
				RegulatorURL: j.RegulatorURL,
				Note:         "Verify license directly with the regulator. v1 does not query regulator APIs.",
				Code:         ErrNotImplemented,
			}
			return okCT(r.RegulatorURL), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "lookup_visa_requirements", Description: "Returns the official citation URL for destination-country visa requirements."},
		func(_ context.Context, _ *mcp.CallToolRequest, a visaArgs) (*mcp.CallToolResult, visaResult, error) {
			j := p.find(a.Destination)
			if j == nil {
				return errCT(ErrJurisdictionUnknown, a.Destination), visaResult{Code: ErrJurisdictionUnknown}, nil
			}
			r := visaResult{CitationURL: j.CitationURL, Note: "v1: link to regulator only."}
			return okCT(r.CitationURL), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "_meta_disclaimer", Description: "Returns the mandatory legal disclaimer the UI must surface."},
		func(_ context.Context, _ *mcp.CallToolRequest, _ disclaimerArgs) (*mcp.CallToolResult, disclaimerResult, error) {
			r := disclaimerResult{Text: p.Disclaimer}
			return okCT(r.Text), r, nil
		})

	return s
}

func (p *pack) find(country string) *Jurisdiction {
	c := strings.ToUpper(country)
	for i := range p.Jurisdictions {
		if strings.EqualFold(p.Jurisdictions[i].CountryCode, c) {
			return &p.Jurisdictions[i]
		}
	}
	return nil
}

func runStdio(s *mcp.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	session, err := s.Connect(ctx, &mcp.StdioTransport{}, nil)
	if err != nil {
		return err
	}
	session.Wait()
	return nil
}

func runHTTP(s *mcp.Server, addr string) error {
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return s }, nil)
	srv := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("mcp-law listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(sctx)
}

func okCT(t string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: t}}}
}
func errCT(code, msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("[%s] %s", code, msg)}}}
}
func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
