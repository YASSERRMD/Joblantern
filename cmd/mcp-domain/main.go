// Package main is the joblantern.domain MCP server.
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
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/domain"
)

const ServerName = "joblantern.domain"

const (
	ErrInvalidDomain = "INVALID_DOMAIN"
	ErrWhois         = "WHOIS_UNAVAILABLE"
	ErrUpstream      = "UPSTREAM_ERROR"
	ErrRateLimited   = "RATE_LIMITED"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8085"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	comp := domain.NewComposer()
	s := newServer(comp)

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

type ageArgs struct {
	Domain string `json:"domain"`
}
type ageResult struct {
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at,omitempty"`
	AgeDays   int    `json:"age_days"`
	Registrar string `json:"registrar,omitempty"`
	Country   string `json:"country_code,omitempty"`
	Code      string `json:"code,omitempty"`
}

type sslArgs struct {
	Domain string `json:"domain"`
}
type sslResult struct {
	*domain.CertSummary
	Code string `json:"code,omitempty"`
}

type archiveArgs struct {
	Domain string `json:"domain"`
}
type archiveResult struct {
	*domain.ArchiveSummary
	Code string `json:"code,omitempty"`
}

type profileArgs struct {
	Domain string `json:"domain"`
}
type profileResult struct {
	*domain.Profile
	Code string `json:"code,omitempty"`
}

func newServer(c *domain.Composer) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "domain_age", Description: "Return creation date, age in days, registrar, registrant country (from WHOIS)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a ageArgs) (*mcp.CallToolResult, ageResult, error) {
			if a.Domain == "" {
				return errCT(ErrInvalidDomain, "domain required"), ageResult{Code: ErrInvalidDomain}, nil
			}
			w, err := c.WHOIS.Lookup(ctx, a.Domain)
			if err != nil {
				if errors.Is(err, domain.ErrInvalidDomain) {
					return errCT(ErrInvalidDomain, err.Error()), ageResult{Code: ErrInvalidDomain}, nil
				}
				return errCT(ErrWhois, err.Error()), ageResult{Code: ErrWhois}, nil
			}
			ageDays := -1
			created := ""
			if !w.CreatedAt.IsZero() {
				ageDays = int(time.Since(w.CreatedAt).Hours() / 24)
				created = w.CreatedAt.UTC().Format(time.RFC3339)
			}
			r := ageResult{Domain: w.Domain, CreatedAt: created, AgeDays: ageDays, Registrar: w.Registrar, Country: w.CountryCode}
			b, _ := json.Marshal(r)
			return okCT(string(b)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "ssl_history", Description: "Return cert count, first/last issuance, unique issuers (crt.sh)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a sslArgs) (*mcp.CallToolResult, sslResult, error) {
			if a.Domain == "" {
				return errCT(ErrInvalidDomain, "domain required"), sslResult{Code: ErrInvalidDomain}, nil
			}
			cs, err := c.CrtSH.Summary(ctx, a.Domain)
			if err != nil {
				return errCT(ErrUpstream, err.Error()), sslResult{Code: ErrUpstream}, nil
			}
			return okCT(fmt.Sprintf("%d certs", cs.CertCount)), sslResult{CertSummary: cs}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "archive_history", Description: "Return earliest/latest Wayback snapshot and snapshot count."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a archiveArgs) (*mcp.CallToolResult, archiveResult, error) {
			if a.Domain == "" {
				return errCT(ErrInvalidDomain, "domain required"), archiveResult{Code: ErrInvalidDomain}, nil
			}
			as, err := c.Wayback.Summary(ctx, a.Domain)
			if err != nil {
				return errCT(ErrUpstream, err.Error()), archiveResult{Code: ErrUpstream}, nil
			}
			return okCT(fmt.Sprintf("%d snapshots", as.SnapshotCount)), archiveResult{ArchiveSummary: as}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "full_domain_profile", Description: "Combined WHOIS + SSL + Wayback with derived age_days and freshness_score."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a profileArgs) (*mcp.CallToolResult, profileResult, error) {
			if a.Domain == "" {
				return errCT(ErrInvalidDomain, "domain required"), profileResult{Code: ErrInvalidDomain}, nil
			}
			p, err := c.FullProfile(ctx, a.Domain)
			if err != nil && p == nil {
				return errCT(ErrUpstream, err.Error()), profileResult{Code: ErrUpstream}, nil
			}
			b, _ := json.Marshal(p)
			return okCT(string(b)), profileResult{Profile: p}, nil
		})

	return s
}

func runStdio(s *mcp.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	session, err := s.Connect(ctx, &mcp.StdioTransport{}, nil)
	if err != nil {
		return err
	}
	_ = session.Wait()
	return nil
}

func runHTTP(s *mcp.Server, addr string) error {
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return s }, nil)
	srv := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("mcp-domain listening", "addr", addr)
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
