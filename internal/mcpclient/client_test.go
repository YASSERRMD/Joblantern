package mcpclient_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/mcpclient"
)

// startFakeServer spins up an in-process MCP server exposing a single
// "echo" tool and returns a wired Client.
type echoArgs struct {
	Message string `json:"message"`
}

type echoResult struct {
	Echo string `json:"echo"`
}

func startFakeServer(t *testing.T) (*mcpclient.Client, func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "fake", Version: "0.0.0"}, nil)
	mcp.AddTool(server,
		&mcp.Tool{Name: "echo", Description: "echo input"},
		func(_ context.Context, _ *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, echoResult, error) {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: args.Message}},
			}, echoResult{Echo: args.Message}, nil
		},
	)

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}

	client, err := mcpclient.Dial(ctx, mcpclient.Config{
		Server:         "fake",
		CallTimeout:    2 * time.Second,
		MaxRetries:     1,
		InitialBackoff: 10 * time.Millisecond,
	}, clientTransport)
	if err != nil {
		_ = serverSession.Close()
		cancel()
		t.Fatalf("client dial: %v", err)
	}

	return client, func() {
		_ = client.Close()
		_ = serverSession.Close()
		serverSession.Wait()
		cancel()
	}
}

func TestCallTool_OK(t *testing.T) {
	client, cleanup := startFakeServer(t)
	defer cleanup()

	var got echoResult
	_, err := client.CallTool(context.Background(), "echo", echoArgs{Message: "hi"}, &got)
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if got.Echo != "hi" {
		t.Errorf("got Echo=%q want %q", got.Echo, "hi")
	}
}

func TestCallTool_UnknownTool(t *testing.T) {
	client, cleanup := startFakeServer(t)
	defer cleanup()

	_, err := client.CallTool(context.Background(), "does-not-exist", echoArgs{Message: "x"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestCallTool_AuditHookFires(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	server := mcp.NewServer(&mcp.Implementation{Name: "fake", Version: "0.0.0"}, nil)
	mcp.AddTool(server,
		&mcp.Tool{Name: "echo"},
		func(_ context.Context, _ *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, echoResult, error) {
			return &mcp.CallToolResult{}, echoResult{Echo: args.Message}, nil
		},
	)
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	var (
		mu      sync.Mutex
		entries []mcpclient.AuditEntry
	)
	client, err := mcpclient.Dial(ctx, mcpclient.Config{
		Server:      "fake",
		CallTimeout: time.Second,
		AuditHook: func(_ context.Context, e mcpclient.AuditEntry) {
			mu.Lock()
			defer mu.Unlock()
			entries = append(entries, e)
		},
	}, clientTransport)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = client.Close() }()

	for i := 0; i < 3; i++ {
		_, _ = client.CallTool(ctx, "echo", echoArgs{Message: "x"}, nil)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(entries) != 3 {
		t.Fatalf("expected 3 audit entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Server != "fake" || e.Tool != "echo" || e.Status != "ok" {
			t.Errorf("entry %d unexpected: %+v", i, e)
		}
		if e.ArgsHash == "" {
			t.Errorf("entry %d missing args hash", i)
		}
		if e.LatencyMS < 0 {
			t.Errorf("entry %d negative latency", i)
		}
	}
}
