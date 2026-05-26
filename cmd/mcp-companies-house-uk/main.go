// Package main is the joblantern.registry.uk MCP server — a
// Companies-House-only variant of mcp-registry. Bundling the
// provider as its own MCP binary lets operators in jurisdictions
// where Companies House is the canonical registry deploy just this
// server without OpenCorporates rate-limit overhead.
//
// The tool surface mirrors mcp-registry exactly so the agent's
// registry sub-agent does not need to special-case it.
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
	"github.com/yasserrmd/joblantern/internal/registry/companieshouse"
)

const ServerName = "joblantern.registry.uk"

const (
	ErrCompanyNotFound = "COMPANY_NOT_FOUND"
	ErrRateLimited     = "RATE_LIMITED"
	ErrTokenInvalid    = "TOKEN_INVALID"
	ErrUpstream        = "UPSTREAM_ERROR"
	ErrInvalidArgs     = "INVALID_ARGS"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8091"), "HTTP listen addr")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	prov := companieshouse.New(os.Getenv("COMPANIES_HOUSE_API_KEY"))
	s := newServer(prov)
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

type lookupArgs struct {
	Name  string `json:"name"`
	Limit int    `json:"limit,omitempty"`
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

func newServer(prov registry.Provider) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "lookup_company", Description: "Search Companies House by name."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a lookupArgs) (*mcp.CallToolResult, lookupResult, error) {
			if a.Name == "" {
				return errCT(ErrInvalidArgs, "name required"), lookupResult{Code: ErrInvalidArgs}, nil
			}
			ms, err := prov.LookupByName(ctx, a.Name, "gb", a.Limit)
			if err != nil {
				return mapLookupErr(err)
			}
			b, _ := json.Marshal(ms)
			return okCT(string(b)), lookupResult{Matches: ms}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "get_company", Description: "Fetch the full Companies House record for an id like \"gb/12345\"."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a getArgs) (*mcp.CallToolResult, getResult, error) {
			if a.ID == "" {
				return errCT(ErrInvalidArgs, "id required"), getResult{Code: ErrInvalidArgs}, nil
			}
			c, err := prov.Get(ctx, a.ID)
			if err != nil {
				return mapGetErr(err)
			}
			return okCT(c.Name), getResult{Company: c}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "check_registration_status", Description: "Derive is_active, is_recent, age_days for a UK company."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a statusArgs) (*mcp.CallToolResult, statusResult, error) {
			c, err := prov.Get(ctx, a.ID)
			if err != nil {
				return mapStatusErr(err)
			}
			active := strings.EqualFold(c.Status, "active")
			age := -1
			if !c.IncorporationDate.IsZero() {
				age = int(time.Since(c.IncorporationDate).Hours() / 24)
			}
			r := statusResult{IsActive: active, IsRecent: age >= 0 && age < 365, AgeDays: age}
			return okCT(fmt.Sprintf("%+v", r)), r, nil
		})

	return s
}

func mapLookupErr(err error) (*mcp.CallToolResult, lookupResult, error) {
	switch {
	case errors.Is(err, registry.ErrNotFound):
		return okCT("no match"), lookupResult{Code: ErrCompanyNotFound}, nil
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
	default:
		return errCT(ErrUpstream, err.Error()), getResult{Code: ErrUpstream}, nil
	}
}

func mapStatusErr(err error) (*mcp.CallToolResult, statusResult, error) {
	switch {
	case errors.Is(err, registry.ErrNotFound):
		return okCT("not found"), statusResult{Code: ErrCompanyNotFound, AgeDays: -1}, nil
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
		slog.Info("mcp-companies-house-uk listening", "addr", addr)
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
