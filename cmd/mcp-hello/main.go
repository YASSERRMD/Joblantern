// Package main is the joblantern.hello MCP server — a single-tool
// reference implementation that exercises both stdio and streamable
// HTTP transports. It exists so that:
//
//  1. New contributors have a complete, runnable example to copy when
//     building real MCP servers (see docs/MCP-SPECS/_TEMPLATE.md).
//  2. Integration tests can dial a real MCP child process without
//     standing up a full upstream dependency.
//
// Usage:
//
//	mcp-hello                       # stdio (default; for spawning from an agent)
//	mcp-hello -transport=http -addr=:8081
package main

import (
	"context"
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
)

// HelloArgs is the structured argument schema for the hello tool.
type HelloArgs struct {
	Name string `json:"name" jsonschema:"the name to greet"`
}

// HelloResult is the structured-content result schema.
type HelloResult struct {
	Greeting string `json:"greeting"`
}

func hello(_ context.Context, _ *mcp.CallToolRequest, args HelloArgs) (*mcp.CallToolResult, HelloResult, error) {
	name := args.Name
	if name == "" {
		name = "world"
	}
	greet := "Hello, " + name + "!"
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: greet}},
	}, HelloResult{Greeting: greet}, nil
}

func newServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "joblantern.hello",
		Version: "0.0.1",
	}, nil)
	mcp.AddTool(s,
		&mcp.Tool{
			Name:        "hello",
			Description: "Greets the supplied name. Returns {greeting: string}.",
		},
		hello,
	)
	return s
}

func main() {
	transport := flag.String("transport", "stdio", "transport: stdio | http")
	addr := flag.String("addr", ":8081", "HTTP listen address when -transport=http")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	switch *transport {
	case "stdio":
		if err := runStdio(); err != nil {
			logger.Error("mcp-hello stdio exit", "err", err)
			os.Exit(1)
		}
	case "http":
		if err := runHTTP(*addr); err != nil {
			logger.Error("mcp-hello http exit", "err", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown transport %q\n", *transport)
		os.Exit(2)
	}
}

func runStdio() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := newServer()
	session, err := server.Connect(ctx, &mcp.StdioTransport{}, nil)
	if err != nil {
		return err
	}
	session.Wait()
	return nil
}

func runHTTP(addr string) error {
	server := newServer()
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0, // streaming
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("mcp-hello listening", "addr", addr)
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
