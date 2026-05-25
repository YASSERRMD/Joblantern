// Package main is the joblantern.address MCP server.
//
// Tools:
//
//	verify_address_exists      forward geocode via Nominatim
//	reverse_geocode            reverse geocode via Nominatim
//	classify_building_type     residential / commercial / mixed via Overpass
//	address_cluster_check      (stub in v1; wired to Postgres in Phase 08)
//	_meta/attribution          required OSM attribution string
//
// Environment:
//
//	NOMINATIM_URL  (default http://localhost:8088)
//	OVERPASS_URL   (default http://localhost:8089/api/interpreter)
//	TRANSPORT      "stdio" (default) | "http"
//	ADDR           HTTP listen address (default :8082)
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/cache"
	"github.com/yasserrmd/joblantern/internal/nominatim"
	"github.com/yasserrmd/joblantern/internal/overpass"
)

// ServerName is the canonical MCP server id.
const ServerName = "joblantern.address"

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8082"), "HTTP listen address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	deps := buildDeps()
	s := newServer(deps)

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
	default:
		logger.Error("unknown transport", "transport", *transport)
		os.Exit(2)
	}
}

type deps struct {
	nom   *nominatim.Client
	over  *overpass.Client
	cache *cache.TTL[string, any]
}

func buildDeps() deps {
	return deps{
		nom:   nominatim.New(getenv("NOMINATIM_URL", "http://localhost:8088")),
		over:  overpass.New(getenv("OVERPASS_URL", "http://localhost:8089/api/interpreter")),
		cache: cache.New[string, any](6 * time.Hour),
	}
}

func newServer(d deps) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)
	addTools(s, d)
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
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("mcp-address listening", "addr", addr)
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
