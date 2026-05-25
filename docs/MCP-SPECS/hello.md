# MCP Server — `joblantern.hello`

## Purpose

Reference implementation. Used as the canonical example when adding a
new MCP server and as a target for end-to-end transport tests.

## Identifier

- **Server name:** `joblantern.hello`
- **Binary:** `cmd/mcp-hello/`
- **Transport (dev):** `stdio`
- **Transport (prod):** `streamable HTTP` on `:8081`

## Tools

### `hello`

- **Description.** Returns a greeting for the supplied name.
- **Args schema.**
  ```jsonc
  { "type": "object", "properties": { "name": { "type": "string" } } }
  ```
- **Result schema.**
  ```jsonc
  { "type": "object", "properties": { "greeting": { "type": "string" } } }
  ```
- **Example.** `{ "name": "joblantern" }` → `{ "greeting": "Hello, joblantern!" }`.
- **Error codes.** None.

## Rate limits

None — local-only reference server.

## License / attribution

Apache 2.0, no upstream attribution required.

## Observability

OTel spans + audit rows on the consumer side via `internal/mcpclient`.

## Citations

—
