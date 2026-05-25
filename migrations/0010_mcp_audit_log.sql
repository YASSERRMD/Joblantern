-- +goose Up
-- mcp_audit_log
-- -------------
-- Append-only log of every MCP tool call made by the agent or by any
-- internal subsystem. This is the audit trail that lets us:
--
--   * Reproduce a verdict — replay the exact tool calls in order.
--   * Diagnose latency outliers per server / tool.
--   * Detect upstream API regressions (jumps in error rate).
--   * Bill / rate-limit by client over time.
--
-- args_hash is a SHA-256 of the canonicalised arguments so we can group
-- by "same call, different verifications" without storing potentially
-- sensitive arguments in plain text. Servers that need to log full
-- arguments can do so in their own table; this one stays a thin index.

CREATE TABLE mcp_audit_log (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id uuid        REFERENCES verifications(id) ON DELETE SET NULL,

    server          text        NOT NULL,
    tool            text        NOT NULL,
    args_hash       text        NOT NULL,

    latency_ms      integer     NOT NULL CHECK (latency_ms >= 0),
    status          text        NOT NULL CHECK (status IN ('ok','error','timeout','rate_limited')),
    error           text,

    called_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX mcp_audit_log_verification_idx ON mcp_audit_log (verification_id);
CREATE INDEX mcp_audit_log_called_at_idx    ON mcp_audit_log (called_at DESC);
CREATE INDEX mcp_audit_log_server_tool_idx  ON mcp_audit_log (server, tool);

-- +goose Down
DROP TABLE IF EXISTS mcp_audit_log;
