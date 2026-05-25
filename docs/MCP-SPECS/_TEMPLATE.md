# MCP Server — `<server-name>`

> Copy this file when adding a new MCP server. Fill in every section.

## Purpose

One paragraph describing what this server does, what upstream data
source(s) it talks to, and which sub-agent consumes it.

## Identifier

- **Server name:** `joblantern.<short>`
- **Binary:** `cmd/mcp-<short>/`
- **Transport (dev):** `stdio`
- **Transport (prod):** `streamable HTTP` on `:<port>`

## Tools

For each tool:

### `tool_name`

- **Description.** One-line summary.
- **Args schema (JSON).**
  ```jsonc
  {
    "type": "object",
    "properties": {
      "field": { "type": "string" }
    },
    "required": ["field"]
  }
  ```
- **Result schema (JSON).**
  ```jsonc
  { "type": "object", "properties": { "ok": { "type": "boolean" } } }
  ```
- **Example.** Input → output.
- **Error codes.** Map of structured codes this tool may return.

## Rate limits

Hard upstream limits (free tier, paid tier, terms-of-use) and how this
server respects them (cache TTL, token-bucket, etc.).

## License / attribution

Upstream license, the attribution string Joblantern must display in the
UI when consuming this server, and any share-alike obligations.

## Observability

- OpenTelemetry spans per tool call.
- Audit-log rows written to `mcp_audit_log` via the shared mcpclient
  wrapper.

## Citations

Authoritative sources backing any factual claim the server can make.
