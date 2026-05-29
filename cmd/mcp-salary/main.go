// Package main is the joblantern.salary MCP server.
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

const ServerName = "joblantern.salary"

const (
	ErrUnknownRole     = "UNKNOWN_ROLE"
	ErrUnknownCountry  = "UNKNOWN_COUNTRY"
	ErrUnknownCurrency = "CURRENCY_UNAVAILABLE"
)

//go:embed bands.json
var defaultBandsJSON []byte

// Pack is the parsed bands file.
type Pack struct {
	SourceNote string             `json:"source_note"`
	Bands      []Band             `json:"bands"`
	FXToUSD    map[string]float64 `json:"fx_to_usd"`
}

// Band is one (country, role) entry.
type Band struct {
	Country   string  `json:"country"`
	Role      string  `json:"role"`
	Currency  string  `json:"currency"`
	P10       float64 `json:"p10"`
	P50       float64 `json:"p50"`
	P90       float64 `json:"p90"`
	SourceURL string  `json:"source_url"`
}

func loadPack() (*Pack, error) {
	var p Pack
	if err := json.Unmarshal(defaultBandsJSON, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8088"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	pk, err := loadPack()
	if err != nil {
		logger.Error("load bands", "err", err)
		os.Exit(1)
	}

	s := newServer(pk)
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

type rangeArgs struct {
	Country       string  `json:"country"`
	Role          string  `json:"role"`
	Currency      string  `json:"currency"`
	ClaimedAmount float64 `json:"claimed_amount"`
	ClaimedPeriod string  `json:"claimed_period,omitempty"` // "month" (default) or "year"
}

type rangeResult struct {
	WithinRange     bool    `json:"within_range"`
	Percentile      float64 `json:"percentile"`
	MultiplierVsP50 float64 `json:"multiplier_vs_p50"`
	SourceURL       string  `json:"source,omitempty"`
	Code            string  `json:"code,omitempty"`
}

type fxArgs struct {
	Amount       float64 `json:"amount"`
	FromCurrency string  `json:"from"`
	ToCurrency   string  `json:"to,omitempty"` // defaults to USD
}

type fxResult struct {
	Amount float64 `json:"amount"`
	Code   string  `json:"code,omitempty"`
}

func newServer(p *Pack) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "salary_range_check", Description: "Compare a claimed salary to the bundled bands for (country, role, currency)."},
		func(_ context.Context, _ *mcp.CallToolRequest, a rangeArgs) (*mcp.CallToolResult, rangeResult, error) {
			b := p.find(a.Country, a.Role)
			if b == nil {
				code := ErrUnknownRole
				if !p.hasCountry(a.Country) {
					code = ErrUnknownCountry
				}
				return errCT(code, "no band"), rangeResult{Code: code}, nil
			}
			amount := a.ClaimedAmount
			if strings.EqualFold(a.ClaimedPeriod, "year") {
				amount /= 12
			}
			// Convert to band currency if needed.
			if !strings.EqualFold(a.Currency, b.Currency) {
				to, ok1 := p.FXToUSD[strings.ToUpper(b.Currency)]
				from, ok2 := p.FXToUSD[strings.ToUpper(a.Currency)]
				if !ok1 || !ok2 || from == 0 || to == 0 {
					return errCT(ErrUnknownCurrency, "no FX"), rangeResult{Code: ErrUnknownCurrency}, nil
				}
				amount = amount * from / to
			}
			multi := amount / b.P50
			percentile := percentile(amount, b)
			within := amount >= b.P10 && amount <= b.P90*1.2 // allow 20% headroom on high
			r := rangeResult{
				WithinRange: within, Percentile: percentile,
				MultiplierVsP50: multi, SourceURL: b.SourceURL,
			}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "currency_normalize", Description: "Convert between bundled currencies (target defaults to USD)."},
		func(_ context.Context, _ *mcp.CallToolRequest, a fxArgs) (*mcp.CallToolResult, fxResult, error) {
			from, ok1 := p.FXToUSD[strings.ToUpper(a.FromCurrency)]
			to := 1.0
			if a.ToCurrency != "" && !strings.EqualFold(a.ToCurrency, "USD") {
				t, ok := p.FXToUSD[strings.ToUpper(a.ToCurrency)]
				if !ok {
					return errCT(ErrUnknownCurrency, "to"), fxResult{Code: ErrUnknownCurrency}, nil
				}
				to = t
			}
			if !ok1 {
				return errCT(ErrUnknownCurrency, "from"), fxResult{Code: ErrUnknownCurrency}, nil
			}
			out := fxResult{Amount: a.Amount * from / to}
			return okCT(fmt.Sprintf("%.2f", out.Amount)), out, nil
		})

	return s
}

func percentile(amount float64, b *Band) float64 {
	switch {
	case amount <= b.P10:
		return 0.1 * amount / b.P10
	case amount <= b.P50:
		return 0.1 + 0.4*(amount-b.P10)/(b.P50-b.P10)
	case amount <= b.P90:
		return 0.5 + 0.4*(amount-b.P50)/(b.P90-b.P50)
	default:
		return 0.9 + 0.1*(amount-b.P90)/(b.P90)
	}
}

func (p *Pack) find(country, role string) *Band {
	for i := range p.Bands {
		b := &p.Bands[i]
		if strings.EqualFold(b.Country, country) && strings.EqualFold(b.Role, role) {
			return b
		}
	}
	return nil
}

func (p *Pack) hasCountry(country string) bool {
	for i := range p.Bands {
		if strings.EqualFold(p.Bands[i].Country, country) {
			return true
		}
	}
	return false
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
		slog.Info("mcp-salary listening", "addr", addr)
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
