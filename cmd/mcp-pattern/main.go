// Package main is the joblantern.pattern MCP server.
//
// Tools:
//
//	analyze_listing_text     scores arbitrary text against the rule pack
//	detect_red_flag_phrases  returns matched phrases without scoring
//	language_mismatch_check  flags Cyrillic/CJK in Latin-script jurisdictions
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

	"github.com/yasserrmd/joblantern/internal/pattern"
)

const ServerName = "joblantern.pattern"

const (
	ErrInvalidArgs      = "INVALID_ARGS"
	ErrUnsupportedLang  = "UNSUPPORTED_LANGUAGE"
	ErrEmbedUnavailable = "EMBEDDING_UNAVAILABLE"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8087"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	rp, err := pattern.DefaultPack()
	if err != nil {
		logger.Error("load rule pack", "err", err)
		os.Exit(1)
	}

	s := newServer(rp)
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

type analyzeArgs struct {
	Text     string `json:"text"`
	Language string `json:"language,omitempty"`
}

type detectArgs struct {
	Text string `json:"text"`
}

type langArgs struct {
	Text         string `json:"text"`
	Jurisdiction string `json:"jurisdiction"`
}

type langResult struct {
	Mismatch bool   `json:"mismatch"`
	Kind     string `json:"kind,omitempty"`
}

func newServer(rp *pattern.RulePack) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "analyze_listing_text", Description: "Score arbitrary listing text for known scam patterns."},
		func(_ context.Context, _ *mcp.CallToolRequest, a analyzeArgs) (*mcp.CallToolResult, pattern.Result, error) {
			if a.Text == "" {
				return errCT(ErrInvalidArgs, "text required"), pattern.Result{}, nil
			}
			r := rp.Analyse(a.Text)
			b, _ := json.Marshal(r)
			return okCT(string(b)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "detect_red_flag_phrases", Description: "Return matched phrases without composite scoring."},
		func(_ context.Context, _ *mcp.CallToolRequest, a detectArgs) (*mcp.CallToolResult, pattern.Result, error) {
			if a.Text == "" {
				return errCT(ErrInvalidArgs, "text required"), pattern.Result{}, nil
			}
			r := rp.Analyse(a.Text)
			// Reuse Result struct but blank the composite score.
			r.CompositeScore = 0
			return okCT(fmt.Sprintf("%d hits", len(r.RedFlags))), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "language_mismatch_check", Description: "Flag Cyrillic/CJK characters in Latin-script jurisdictions."},
		func(_ context.Context, _ *mcp.CallToolRequest, a langArgs) (*mcp.CallToolResult, langResult, error) {
			if a.Text == "" {
				return errCT(ErrInvalidArgs, "text required"), langResult{}, nil
			}
			m, k := pattern.LanguageMismatchCheck(a.Text, a.Jurisdiction)
			return okCT(fmt.Sprintf("%v %s", m, k)), langResult{Mismatch: m, Kind: k}, nil
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
		slog.Info("mcp-pattern listening", "addr", addr)
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
