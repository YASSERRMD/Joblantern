// Package main is the joblantern.routing MCP server.
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

	"github.com/yasserrmd/joblantern/internal/ors"
)

const ServerName = "joblantern.routing"

const (
	ErrRateLimited  = "RATE_LIMITED"
	ErrOutOfRegion  = "OUT_OF_REGION"
	ErrInvalidCoord = "INVALID_COORD"
	ErrUpstream     = "UPSTREAM_ERROR"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8090"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cli := ors.New(os.Getenv("ORS_API_KEY"))
	if cli.APIKey == "" {
		logger.Warn("ORS_API_KEY not set; route tool will return UPSTREAM_ERROR")
	}

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

type point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type routeArgs struct {
	From point  `json:"from"`
	To   point  `json:"to"`
	Mode string `json:"mode,omitempty"` // driving | walking | cycling
}

type routeResult struct {
	DistanceKM float64 `json:"distance_km"`
	DurationM  float64 `json:"duration_min"`
	Code       string  `json:"code,omitempty"`
}

type commuteArgs struct {
	From point  `json:"from"`
	To   point  `json:"to"`
	Mode string `json:"mode,omitempty"`
}

type commuteResult struct {
	Reachable  bool    `json:"reachable"`
	Plausible  bool    `json:"plausible"`
	DurationM  float64 `json:"duration_min"`
	DistanceKM float64 `json:"distance_km"`
	Code       string  `json:"code,omitempty"`
}

func parseMode(s string) ors.Mode {
	switch strings.ToLower(s) {
	case "walking", "foot":
		return ors.ModeWalking
	case "cycling", "bike":
		return ors.ModeCycling
	default:
		return ors.ModeDriving
	}
}

func newServer(cli *ors.Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "route", Description: "Distance and duration between two coordinates (driving/walking/cycling)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a routeArgs) (*mcp.CallToolResult, routeResult, error) {
			r, err := cli.Route(ctx, parseMode(a.Mode), a.From.Lat, a.From.Lon, a.To.Lat, a.To.Lon)
			if err != nil {
				return mapErrRoute(err)
			}
			out := routeResult{DistanceKM: r.DistanceM / 1000.0, DurationM: r.DurationS / 60.0}
			return okCT(jstr(out)), out, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "commute_realism_check", Description: "Derived: is the commute from home to office reachable and < 3 hours one-way?"},
		func(ctx context.Context, _ *mcp.CallToolRequest, a commuteArgs) (*mcp.CallToolResult, commuteResult, error) {
			r, err := cli.Route(ctx, parseMode(a.Mode), a.From.Lat, a.From.Lon, a.To.Lat, a.To.Lon)
			if err != nil {
				_, out, _ := mapErrCommute(err)
				return errCTCommute(err), out, nil
			}
			durMin := r.DurationS / 60.0
			out := commuteResult{
				Reachable:  true,
				Plausible:  durMin < 180,
				DurationM:  durMin,
				DistanceKM: r.DistanceM / 1000.0,
			}
			return okCT(jstr(out)), out, nil
		})

	return s
}

func mapErrRoute(err error) (*mcp.CallToolResult, routeResult, error) {
	switch {
	case errors.Is(err, ors.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), routeResult{Code: ErrRateLimited}, nil
	case errors.Is(err, ors.ErrOutOfRegion):
		return errCT(ErrOutOfRegion, err.Error()), routeResult{Code: ErrOutOfRegion}, nil
	default:
		return errCT(ErrUpstream, err.Error()), routeResult{Code: ErrUpstream}, nil
	}
}

func mapErrCommute(err error) (*mcp.CallToolResult, commuteResult, error) {
	switch {
	case errors.Is(err, ors.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), commuteResult{Code: ErrRateLimited}, nil
	case errors.Is(err, ors.ErrOutOfRegion):
		return errCT(ErrOutOfRegion, err.Error()), commuteResult{Code: ErrOutOfRegion}, nil
	default:
		return errCT(ErrUpstream, err.Error()), commuteResult{Code: ErrUpstream}, nil
	}
}

func errCTCommute(err error) *mcp.CallToolResult { return errCT(ErrUpstream, err.Error()) }

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
		slog.Info("mcp-routing listening", "addr", addr)
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
func jstr(v any) string { b, _ := json.Marshal(v); return string(b) }
func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
