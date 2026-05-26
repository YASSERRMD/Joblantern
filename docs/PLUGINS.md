# Joblantern plugins

Third parties can contribute MCP servers without forking Joblantern by
publishing a small YAML manifest plus the server binary. Operators
register plugins per-deployment.

## Authoring a plugin

1. Build an MCP server using the official Go SDK (or any other MCP
   implementation). It must speak either `stdio` or streamable HTTP.
2. Write a `joblantern-mcp.yaml` manifest. See
   [`examples/community-mcp/joblantern-mcp.yaml`](../examples/community-mcp/joblantern-mcp.yaml).
3. License the server under a permissive SPDX id (`Apache-2.0`, `MIT`,
   `BSD-2-Clause`, `BSD-3-Clause`, `ISC`, or `MPL-2.0`). The manifest
   loader rejects everything else outright.
4. Sign the manifest (recommended): strip the `signature_b64` /
   `pubkey_hex` fields, sign the canonical YAML with an ed25519 key,
   and paste them back.

## Trust levels

| Trust | When to use |
|---|---|
| `official` | Maintained by the Joblantern project. Allowed unsandboxed by default. |
| `community` | Operator has reviewed it. Signature required. |
| `external` | Unreviewed. **Run in a sandbox** (gVisor / Firecracker) and weight evidence low. |

The agent already understands per-plugin weights via the risk engine;
operators tune them through the future plugins admin UI.

## Storage

Each registered plugin is one row in the `plugins` table (migration
0013). The full manifest is stored verbatim so an auditor can re-verify
it later without re-fetching from the source URL.

## Why not just download and run anything?

Because the Joblantern binary holds users' verdicts and a curated scam
DB. A malicious MCP server with unbounded host access could
exfiltrate everything. By requiring permissive licensing, signed
manifests, an explicit trust label, and a sandbox recommendation, we
make the threat model legible to operators.

## ADR

A separate ADR (deferred to v0.3 alongside the admin UI) will lock in
the sandboxing policy and the plugin curation process.
