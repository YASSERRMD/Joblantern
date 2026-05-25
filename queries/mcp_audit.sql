-- Append-only writes plus a few read paths used by the observability
-- subsystem (Phase 18) and the replay tooling.

-- name: InsertMCPAuditLog :exec
INSERT INTO mcp_audit_log (
    verification_id,
    server,
    tool,
    args_hash,
    latency_ms,
    status,
    error
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: ListMCPAuditByVerification :many
SELECT *
  FROM mcp_audit_log
 WHERE verification_id = $1
 ORDER BY called_at ASC;

-- name: CountMCPAuditByServerToolStatus :many
-- Used by Prometheus-style introspection panels.
SELECT server,
       tool,
       status,
       COUNT(*)::bigint AS n
  FROM mcp_audit_log
 WHERE called_at >= $1
 GROUP BY server, tool, status
 ORDER BY server, tool, status;
