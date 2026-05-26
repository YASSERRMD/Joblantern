// Package main is the joblantern.vies MCP server — validates EU VAT
// numbers against the official VIES service. One tool only:
//
//	validate_vat_number(country, vat_number) → {valid, name, address}
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

	"github.com/yasserrmd/joblantern/internal/vies"
)

const ServerName = "joblantern.vies"

const (
	ErrInvalidCountry = "INVALID_COUNTRY"
	ErrNotFound       = "VAT_NOT_FOUND"
	ErrUpstream       = "UPSTREAM_ERROR"
	ErrInvalidArgs    = "INVALID_ARGS"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8092"), "HTTP listen addr")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cli := vies.New()
	s := newServer(cli)
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

type args struct {
	Country   string `json:"country"`
	VATNumber string `json:"vat_number"`
}

type result struct {
	*vies.Result
	Code string `json:"code,omitempty"`
}

func newServer(cli *vies.Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)
	mcp.AddTool(s,
		&mcp.Tool{Name: "validate_vat_number", Description: "Validate an EU VAT number (country, vat_number) against the official VIES service."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a args) (*mcp.CallToolResult, result, error) {
			if a.Country == "" || a.VATNumber == "" {
				return errCT(ErrInvalidArgs, "country and vat_number required"), result{Code: ErrInvalidArgs}, nil
			}
			r, err := cli.Validate(ctx, a.Country, a.VATNumber)
			if err != nil {
				switch {
				case errors.Is(err, vies.ErrInvalidCountry):
					return errCT(ErrInvalidCountry, err.Error()), result{Code: ErrInvalidCountry}, nil
				case errors.Is(err, vies.ErrNotFound):
					return errCT(ErrNotFound, err.Error()), result{Code: ErrNotFound}, nil
				default:
					return errCT(ErrUpstream, err.Error()), result{Code: ErrUpstream}, nil
				}
			}
			b, _ := json.Marshal(r)
			return okCT(string(b)), result{Result: r}, nil
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
		slog.Info("mcp-eu-vies listening", "addr", addr)
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
