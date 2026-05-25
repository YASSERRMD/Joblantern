// Package main starts the Joblantern HTTP server.
//
// Phase 01 ships only a healthcheck endpoint and a graceful shutdown
// scaffold. Sub-agents, MCP clients, database access, web UI and the
// rest of the stack are introduced in later phases as documented in
// docs/CLAUDE_CODE_PROMPT_PACK.md.
package main

import (
	"context"
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
)

// healthcheckFlag, when set, performs a one-shot probe against the local
// server and exits 0/1. The Dockerfile.dev HEALTHCHECK uses this flag so
// the distroless image does not need a shell or curl.
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
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func run(addr string, logger *slog.Logger) error {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
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

// probe issues GET /healthz against the local server and returns 0 on OK.
func probe(addr string) int {
	url := "http://127.0.0.1" + addr + "/healthz"
	if len(addr) > 0 && addr[0] != ':' {
		url = "http://" + addr + "/healthz"
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url) //nolint:gosec // local-only probe against own server
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
