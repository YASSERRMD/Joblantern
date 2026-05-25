//go:build integration
// +build integration

package mcpclient_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/mcpclient"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find go.mod walking up from %s", wd)
		}
		dir = parent
	}
}

// TestDialCommand_Hello builds the mcp-hello binary and dials it as a
// child process over stdio. This proves the end-to-end CommandTransport
// path works (separate from the in-memory unit tests).
//
// Skipped when CGO_ENABLED isn't available or on Windows (CI is Linux).
func TestDialCommand_Hello(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build mcp-hello into a temp dir.
	tmp := t.TempDir()
	bin := filepath.Join(tmp, "mcp-hello")
	build := exec.CommandContext(ctx, "go", "build", "-o", bin, "./cmd/mcp-hello")
	build.Dir = repoRoot(t)
	out, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build mcp-hello: %v\n%s", err, out)
	}

	client, err := mcpclient.DialCommand(ctx, mcpclient.Config{
		Server:      "joblantern.hello",
		CallTimeout: 5 * time.Second,
	}, bin)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	type result struct {
		Greeting string `json:"greeting"`
	}
	var got result
	res, err := client.CallTool(ctx, "hello", map[string]any{"name": "joblantern"}, &got)
	if err != nil {
		t.Fatalf("call hello: %v", err)
	}
	if got.Greeting != "Hello, joblantern!" {
		t.Errorf("greeting=%q want %q", got.Greeting, "Hello, joblantern!")
	}
	// Text content should also carry the greeting.
	if len(res.Content) == 0 {
		t.Fatal("no content returned")
	}
	if tc, ok := res.Content[0].(*mcp.TextContent); !ok || tc.Text == "" {
		t.Errorf("unexpected content: %#v", res.Content[0])
	}
}
