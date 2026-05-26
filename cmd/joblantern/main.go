// Package main starts the Joblantern HTTP server.
//
// Phase 13: chi router + JSON API for the agent orchestrator. The
// orchestrator ships with two pure-in-process sub-agents (pattern,
// language) so the binary is useful without external services. MCP-
// backed sub-agents (address, registry, domain, salary, law, routing)
// are wired in their own binaries and dialled by the agent through
// internal/mcpclient.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/yasserrmd/joblantern/internal/agent"
	"github.com/yasserrmd/joblantern/internal/pattern"
	"github.com/yasserrmd/joblantern/internal/risk"
	"github.com/yasserrmd/joblantern/internal/transparency"
	"github.com/yasserrmd/joblantern/internal/web"
)

var healthcheckFlag = flag.Bool("healthcheck", false, "run a one-shot health probe against the local server and exit")

func main() {
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	addr := getenv("JOBLANTERN_ADDR", ":8080")

	if *healthcheckFlag {
		os.Exit(probe(addr))
	}
	if err := run(addr, logger); err != nil {
		logger.Error("server exited", "err", err)
		os.Exit(1)
	}
}

func run(addr string, logger *slog.Logger) error {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(web.SecurityHeaders)
	metrics := web.NewMetrics()
	r.Use(metrics.Middleware)

	limit := web.NewIPRateLimiter(20, time.Minute)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/verify" || r.URL.Path == "/api/v1/verify" {
				limit.Middleware(next).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "ok")
	})
	r.Handle("/metrics", metrics.Endpoint())

	store := agent.NewMemoryStore()
	subs, err := buildBuiltinSubagents()
	if err != nil {
		return err
	}
	orch := agent.New(subs...).WithScorer(func(facts []agent.Fact) (string, float64, []string) {
		o := risk.Score(facts, risk.DefaultBands)
		return o.OverallRisk, o.Confidence, o.Reasons
	})
	api := web.NewAPIHandler(r, store, orch)
	if _, err := web.NewUI(r, store, api); err != nil {
		return err
	}
	if err := web.MountStatic(r); err != nil {
		return err
	}

	// Recruiter signed-badge issuer. v0.1 generates a fresh keypair per
	// process; production should load a persistent ed25519 key from
	// disk (JOBLANTERN_BADGE_KEY, follow-up PR).
	bpub, bpriv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	web.NewBadgeIssuer(r, getenv("JOBLANTERN_BADGE_ISSUER", "http://localhost:8080"), bpriv, bpub, 90*24*time.Hour)

	// Hostile-network surface: /panic-wipe + /lite (no JS, no CSS).
	web.NewHostile(r, store)

	// NGO operator tools: kiosk mode + single-page PDF printout.
	web.NewNGO(r, store, api)

	// Public transparency dashboard — reads anonymised aggregates from
	// the agent store. Future PR replaces the in-memory source with a
	// materialised view fed by a nightly aggregation job.
	web.NewTransparency(r, func() []transparency.Verdict {
		recs, _ := store.List(context.Background(), 0)
		out := make([]transparency.Verdict, 0, len(recs))
		for _, rec := range recs {
			if rec.Verdict == nil {
				continue
			}
			completed := rec.UpdatedAt
			out = append(out, transparency.Verdict{
				CompletedAt: completed,
				Country:     rec.Submission.Jurisdiction,
				Risk:        rec.Verdict.OverallRisk,
			})
		}
		return out
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("joblantern listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	logger.Info("joblantern stopped cleanly")
	return nil
}

func buildBuiltinSubagents() ([]agent.Subagent, error) {
	pp, err := pattern.DefaultPack()
	if err != nil {
		return nil, err
	}
	return []agent.Subagent{
		&agent.MCPSubagent{
			NameStr: "pattern",
			Run_: func(_ context.Context, sub agent.Submission) []agent.Fact {
				if sub.ListingText == "" {
					return nil
				}
				res := pp.Analyse(sub.ListingText)
				if len(res.RedFlags) == 0 {
					return []agent.Fact{{
						Source: "joblantern.pattern", ToolName: "analyze_listing_text",
						FactType: "pattern.red_flag", Value: 0,
						SupportsRisk: "green", Weight: 0.2,
					}}
				}
				out := make([]agent.Fact, 0, len(res.RedFlags))
				for _, h := range res.RedFlags {
					out = append(out, agent.Fact{
						Source: "joblantern.pattern", ToolName: "analyze_listing_text",
						FactType:     "pattern.red_flag",
						Value:        map[string]any{"code": h.Code, "span": h.Span, "description": h.Description},
						SupportsRisk: "red",
						Weight:       h.Weight,
					})
				}
				return out
			},
		},
		&agent.MCPSubagent{
			NameStr: "language",
			Run_: func(_ context.Context, sub agent.Submission) []agent.Fact {
				if sub.ListingText == "" || sub.Jurisdiction == "" {
					return nil
				}
				m, kind := pattern.LanguageMismatchCheck(sub.ListingText, sub.Jurisdiction)
				if !m {
					return nil
				}
				return []agent.Fact{{
					Source: "joblantern.pattern", ToolName: "language_mismatch_check",
					FactType:     "pattern.language_mismatch",
					Value:        map[string]any{"kind": kind, "jurisdiction": sub.Jurisdiction},
					SupportsRisk: "red",
					Weight:       0.5,
				}}
			},
		},
	}, nil
}

func probe(addr string) int {
	url := "http://127.0.0.1" + addr + "/healthz"
	if len(addr) > 0 && addr[0] != ':' {
		url = "http://" + addr + "/healthz"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url) //nolint:gosec
	if err != nil {
		return 1
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
