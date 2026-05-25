// Package main is the joblantern.registry MCP server.
//
// Backed by a registry.Provider; ships with the OpenCorporates
// implementation. Future providers (Companies House, SEC EDGAR…) can
// be plugged in without changing the tool surface.
package main

import (
	"context"
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

	"github.com/yasserrmd/joblantern/internal/registry"
	"github.com/yasserrmd/joblantern/internal/registry/opencorporates"
)

const ServerName = "joblantern.registry"

const (
	ErrCompanyNotFound  = "COMPANY_NOT_FOUND"
	ErrJurisdictionMiss = "JURISDICTION_UNKNOWN"
	ErrRateLimited      = "RATE_LIMITED"
	ErrTokenInvalid     = "TOKEN_INVALID"
	ErrUpstream         = "UPSTREAM_ERROR"
	ErrInvalidArgs      = "INVALID_ARGS"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8084"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	prov := opencorporates.New(os.Getenv("OPENCORPORATES_TOKEN"))

	s := newServer(prov)

	switch *transport {
	case "stdio":
		if err := runStdio(s); err != nil {
			logger.Error("stdio exit", "err", err)
			os.Exit(1)
		}
	case "http":
		if err := runHTTP(s, *addr); err != nil {
			logger.Error("http exit", "err", err)
			os.Exit(1)
		}
	}
}

type lookupArgs struct {
	Name         string `json:"name"`
	Jurisdiction string `json:"jurisdiction,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

type lookupResult struct {
	Matches []registry.Match `json:"matches"`
	Code    string           `json:"code,omitempty"`
}

type getArgs struct {
	ID string `json:"id"`
}

type getResult struct {
	Company *registry.Company `json:"company,omitempty"`
	Code    string            `json:"code,omitempty"`
}

type statusArgs struct {
	ID string `json:"id"`
}

type statusResult struct {
	IsActive bool   `json:"is_active"`
	IsRecent bool   `json:"is_recent"`
	AgeDays  int    `json:"age_days"`
	Code     string `json:"code,omitempty"`
}

type attribArgs struct{}
type attribResult struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func newServer(prov registry.Provider) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "lookup_company", Description: "Search the configured registry by name (optionally scoped to a jurisdiction)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a lookupArgs) (*mcp.CallToolResult, lookupResult, error) {
			if a.Name == "" {
				return errCT(ErrInvalidArgs, "name is required"), lookupResult{Code: ErrInvalidArgs}, nil
			}
			ms, err := prov.LookupByName(ctx, a.Name, a.Jurisdiction, a.Limit)
			if err != nil {
				return mapLookupErr(err)
			}
			b, _ := json.Marshal(ms)
			return okCT(string(b)), lookupResult{Matches: ms}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "get_company", Description: "Fetch a full company record by provider-specific id (e.g. \"gb/12345\")."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a getArgs) (*mcp.CallToolResult, getResult, error) {
			if a.ID == "" {
				return errCT(ErrInvalidArgs, "id is required"), getResult{Code: ErrInvalidArgs}, nil
			}
			c, err := prov.Get(ctx, a.ID)
			if err != nil {
				return mapGetErr(err)
			}
			return okCT(c.Name), getResult{Company: c}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "check_registration_status", Description: "Derive is_active, is_recent (incorporated <12 months) for a company id."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a statusArgs) (*mcp.CallToolResult, statusResult, error) {
			if a.ID == "" {
				return errCT(ErrInvalidArgs, "id is required"), statusResult{Code: ErrInvalidArgs}, nil
			}
			c, err := prov.Get(ctx, a.ID)
			if err != nil {
				return mapStatusErr(err)
			}
			active := strings.EqualFold(c.Status, "active")
			ageDays := -1
			if !c.IncorporationDate.IsZero() {
				ageDays = int(time.Since(c.IncorporationDate).Hours() / 24)
			}
			recent := ageDays >= 0 && ageDays < 365
			r := statusResult{IsActive: active, IsRecent: recent, AgeDays: ageDays}
			b, _ := json.Marshal(r)
			return okCT(string(b)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "_meta_attribution", Description: "Powered-by-OpenCorporates attribution string for any UI displaying these results."},
		func(_ context.Context, _ *mcp.CallToolRequest, _ attribArgs) (*mcp.CallToolResult, attribResult, error) {
			r := attribResult{
				Text: "Powered by OpenCorporates",
				URL:  "https://opencorporates.com/",
			}
			return okCT(r.Text), r, nil
		})

	return s
}

func mapLookupErr(err error) (*mcp.CallToolResult, lookupResult, error) {
	switch {
	case errors.Is(err, registry.ErrNotFound):
		return okCT("no match"), lookupResult{Code: ErrCompanyNotFound}, nil
	case errors.Is(err, registry.ErrJurisdiction):
		return errCT(ErrJurisdictionMiss, err.Error()), lookupResult{Code: ErrJurisdictionMiss}, nil
	case errors.Is(err, registry.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), lookupResult{Code: ErrRateLimited}, nil
	case errors.Is(err, registry.ErrTokenInvalid):
		return errCT(ErrTokenInvalid, err.Error()), lookupResult{Code: ErrTokenInvalid}, nil
	default:
		return errCT(ErrUpstream, err.Error()), lookupResult{Code: ErrUpstream}, nil
	}
}

func mapGetErr(err error) (*mcp.CallToolResult, getResult, error) {
	switch {
	case errors.Is(err, registry.ErrNotFound):
		return okCT("not found"), getResult{Code: ErrCompanyNotFound}, nil
	case errors.Is(err, registry.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), getResult{Code: ErrRateLimited}, nil
	case errors.Is(err, registry.ErrTokenInvalid):
		return errCT(ErrTokenInvalid, err.Error()), getResult{Code: ErrTokenInvalid}, nil
	default:
		return errCT(ErrUpstream, err.Error()), getResult{Code: ErrUpstream}, nil
	}
}

func mapStatusErr(err error) (*mcp.CallToolResult, statusResult, error) {
	switch {
	case errors.Is(err, registry.ErrNotFound):
		return okCT("not found"), statusResult{Code: ErrCompanyNotFound, AgeDays: -1}, nil
	case errors.Is(err, registry.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), statusResult{Code: ErrRateLimited, AgeDays: -1}, nil
	default:
		return errCT(ErrUpstream, err.Error()), statusResult{Code: ErrUpstream, AgeDays: -1}, nil
	}
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
		slog.Info("mcp-registry listening", "addr", addr)
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
