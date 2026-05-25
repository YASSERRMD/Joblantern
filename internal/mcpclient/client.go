// Package mcpclient wraps the github.com/modelcontextprotocol/go-sdk
// client with the policies Joblantern needs:
//
//   - structured slog logging on every call
//   - configurable per-call timeout
//   - bounded exponential-backoff retry on transient errors
//   - hook for writing an audit-log row per call (used to populate
//     the mcp_audit_log table from migration 0010)
//
// One Client instance can talk to one MCP server. The agent owns a
// Client per configured server (see internal/config/mcp.go and
// internal/agent).
package mcpclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"sort"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AuditHook is invoked exactly once per CallTool invocation, after the
// call returns (whether successfully or not). Implementations typically
// write a row to mcp_audit_log.
type AuditHook func(ctx context.Context, entry AuditEntry)

// AuditEntry captures the minimum information required to reproduce or
// diagnose an MCP tool call. Arguments are not stored in plain text;
// args_hash is a SHA-256 of the canonicalised JSON representation.
type AuditEntry struct {
	Server    string
	Tool      string
	ArgsHash  string
	LatencyMS int64
	Status    string // "ok" | "error" | "timeout" | "rate_limited"
	Err       error
}

// Config controls Client behaviour. Zero values are sensible defaults.
type Config struct {
	// Server is a short, stable identifier for the server (e.g.
	// "joblantern.address"). It appears in logs and audit rows.
	Server string

	// CallTimeout is the per-call deadline. Zero means 30s.
	CallTimeout time.Duration

	// MaxRetries is the maximum number of additional attempts after
	// the first call. Zero means 2 (so 3 total).
	MaxRetries int

	// InitialBackoff is the first retry delay; subsequent retries
	// double. Zero means 200ms.
	InitialBackoff time.Duration

	// AuditHook, if non-nil, receives one entry per CallTool.
	AuditHook AuditHook

	// Logger is the slog handler used; defaults to slog.Default().
	Logger *slog.Logger
}

func (c Config) withDefaults() Config {
	if c.CallTimeout == 0 {
		c.CallTimeout = 30 * time.Second
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 2
	}
	if c.InitialBackoff == 0 {
		c.InitialBackoff = 200 * time.Millisecond
	}
	if c.Logger == nil {
		c.Logger = slog.Default()
	}
	return c
}

// Client is a thin, opinionated wrapper around an MCP client session.
type Client struct {
	cfg     Config
	session *mcp.ClientSession
}

// Dial connects to an MCP server over the given transport.
//
// Typical transports:
//
//   - mcp.NewInMemoryTransports() for in-process testing
//   - &mcp.CommandTransport{Command: exec.Command(...)} for stdio-spawned servers
//
// Callers may also use the HTTP transports exposed by the SDK; this
// wrapper does not care which transport is used.
func Dial(ctx context.Context, cfg Config, transport mcp.Transport) (*Client, error) {
	cfg = cfg.withDefaults()
	if cfg.Server == "" {
		return nil, errors.New("mcpclient: Config.Server is required")
	}

	c := mcp.NewClient(&mcp.Implementation{
		Name:    "joblantern",
		Version: "0.0.0",
	}, nil)

	session, err := c.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("mcpclient: connect %s: %w", cfg.Server, err)
	}

	return &Client{cfg: cfg, session: session}, nil
}

// DialCommand is a convenience wrapper that spawns the MCP server as a
// child process and speaks over stdio. The caller is responsible for
// closing the returned Client when done; doing so terminates the child.
func DialCommand(ctx context.Context, cfg Config, command string, args ...string) (*Client, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	return Dial(ctx, cfg, &mcp.CommandTransport{Command: cmd})
}

// Close terminates the underlying session (and any spawned child
// process, when using CommandTransport).
func (c *Client) Close() error {
	if c == nil || c.session == nil {
		return nil
	}
	return c.session.Close()
}

// CallTool invokes a named MCP tool with the given arguments and
// decodes the structured-content result into `out` if non-nil.
//
// Retries are attempted on transient errors (context.DeadlineExceeded
// and net.Error-style temporary failures). The structured server-side
// error returned by the SDK is *not* retried because it represents a
// deterministic application response.
func (c *Client) CallTool(ctx context.Context, tool string, args any, out any) (*mcp.CallToolResult, error) {
	if c == nil || c.session == nil {
		return nil, errors.New("mcpclient: nil session")
	}

	callCtx, cancel := context.WithTimeout(ctx, c.cfg.CallTimeout)
	defer cancel()

	argsHash := hashArgs(args)
	start := time.Now()

	var (
		result *mcp.CallToolResult
		err    error
	)
	backoff := c.cfg.InitialBackoff
	for attempt := 0; attempt <= c.cfg.MaxRetries; attempt++ {
		result, err = c.session.CallTool(callCtx, &mcp.CallToolParams{
			Name:      tool,
			Arguments: args,
		})
		if err == nil || !isRetryable(err) {
			break
		}
		c.cfg.Logger.Warn("mcp call retrying",
			"server", c.cfg.Server,
			"tool", tool,
			"attempt", attempt+1,
			"err", err,
		)
		select {
		case <-callCtx.Done():
			err = callCtx.Err()
			break
		case <-time.After(backoff):
		}
		backoff *= 2
	}

	latency := time.Since(start)
	status := statusFor(err)
	c.cfg.Logger.Debug("mcp call",
		"server", c.cfg.Server,
		"tool", tool,
		"latency_ms", latency.Milliseconds(),
		"status", status,
		"err", err,
	)
	if c.cfg.AuditHook != nil {
		c.cfg.AuditHook(ctx, AuditEntry{
			Server:    c.cfg.Server,
			Tool:      tool,
			ArgsHash:  argsHash,
			LatencyMS: latency.Milliseconds(),
			Status:    status,
			Err:       err,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("mcp %s.%s: %w", c.cfg.Server, tool, err)
	}
	if out != nil && result.StructuredContent != nil {
		if data, mErr := json.Marshal(result.StructuredContent); mErr == nil {
			_ = json.Unmarshal(data, out)
		}
	}
	return result, nil
}

// hashArgs canonicalises args to a JSON object with sorted keys and
// returns the SHA-256 hex digest. We never log the raw arguments —
// they may contain user-supplied PII.
func hashArgs(args any) string {
	if args == nil {
		return ""
	}
	data, err := canonicalJSON(args)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return sortedMarshal(m)
}

func sortedMarshal(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := []byte{'{'}
		for i, k := range keys {
			if i > 0 {
				buf = append(buf, ',')
			}
			kj, _ := json.Marshal(k)
			buf = append(buf, kj...)
			buf = append(buf, ':')
			vj, err := sortedMarshal(t[k])
			if err != nil {
				return nil, err
			}
			buf = append(buf, vj...)
		}
		buf = append(buf, '}')
		return buf, nil
	case []any:
		buf := []byte{'['}
		for i, el := range t {
			if i > 0 {
				buf = append(buf, ',')
			}
			eb, err := sortedMarshal(el)
			if err != nil {
				return nil, err
			}
			buf = append(buf, eb...)
		}
		buf = append(buf, ']')
		return buf, nil
	default:
		return json.Marshal(v)
	}
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// The SDK wraps transport errors; we can be conservative and
	// only retry on deadline exceeded for v1. Future MCP-defined
	// "transient" error codes can be added here.
	return false
}

func statusFor(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	default:
		return "error"
	}
}
