package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/config"
)

func TestLoadMCP_OK(t *testing.T) {
	cfg, err := config.LoadMCP(filepath.Join("..", "..", "config", "mcp.yaml"))
	if err != nil {
		t.Fatalf("LoadMCP: %v", err)
	}
	if len(cfg.Servers) == 0 {
		t.Fatal("expected at least one server")
	}
	got := cfg.EnabledByName("joblantern.hello")
	if got == nil {
		t.Fatal("joblantern.hello not enabled in mcp.yaml")
	}
	if got.Transport != config.TransportStdio {
		t.Errorf("hello transport=%q want stdio", got.Transport)
	}
	if got.CallTimeout != 10*time.Second {
		t.Errorf("hello call_timeout=%v want 10s", got.CallTimeout)
	}
}

func TestLoadMCP_Validation(t *testing.T) {
	bad := t.TempDir()
	path := filepath.Join(bad, "mcp.yaml")
	if err := writeFile(path, "servers:\n  - name: a\n    transport: stdio\n"); err != nil {
		t.Fatal(err)
	}
	_, err := config.LoadMCP(path)
	if err == nil {
		t.Fatal("expected error: stdio transport without command")
	}
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}
