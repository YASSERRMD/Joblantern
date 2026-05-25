# ADR-003 — MCP server pattern: one binary per server

- **Status:** Accepted
- **Date:** 2025-11-09

## Context

Joblantern's agent talks to ~9 specialised MCP servers (address,
streetview, registry, domain, scam-db, salary, law, pattern, routing,
plus the reference hello). Each has a different upstream license, rate
limit, cold-start time, and operational profile. We need to decide
whether they share a process or live as separate binaries.

## Decision

- **One Go binary per MCP server,** living under `cmd/mcp-<name>/`.
- **Stdio transport in development** (the agent spawns each server as a
  child process via `mcp.CommandTransport`).
- **Streamable HTTP transport in production** (each server runs as its
  own container, scaled independently).
- **A shared client wrapper** in `internal/mcpclient` enforces timeout,
  retry, structured logging, and audit-log writes across all callers.
- **Configuration** for the registry lives in `config/mcp.yaml`,
  loaded by `internal/config/mcp.go`.

## Rationale

- **Blast radius.** A misbehaving upstream (rate-limited, OOM, slow)
  cannot take down the rest of the agent.
- **Independent deployment.** Different servers want different release
  cadences and resource limits.
- **License hygiene.** Each binary's own dependency tree is
  license-audited in isolation.
- **Reuse outside Joblantern.** Each MCP server is consumable by Claude
  Desktop, Cursor, other agents — not just our own.

## Consequences

### Positive

- Per-server SLOs.
- Small, focused dependency graphs per binary.
- The hello server demonstrates the pattern end-to-end.

### Negative

- More Docker images, more CI build minutes. Mitigated by sharing the
  multi-stage build cache and using distroless runtime images.
- More process boundaries to debug. Mitigated by the shared mcpclient
  wrapper which standardises logging, tracing, and audit.

## References

- MCP Go SDK: <https://github.com/modelcontextprotocol/go-sdk>
- `docs/MCP-SPECS/_TEMPLATE.md` — spec template every server fills in.
- `docs/MCP-SPECS/hello.md` — reference spec.
