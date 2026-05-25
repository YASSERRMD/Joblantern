# ADR-001 — Stack choices for Joblantern v1

- **Status:** Accepted
- **Date:** 2025-11-09 (locked in the prompt pack)
- **Deciders:** Mohamed Yasser (author)
- **Supersedes:** —
- **Superseded by:** —

## Context

Joblantern is an open-source, Apache 2.0 system whose users include
NGOs, migrant-worker advocacy groups, and individual job seekers who
cannot afford fragile or expensive infrastructure. Every architectural
choice has to optimise for **auditability, low total cost of ownership,
permissive licensing, and operational predictability** — in that order.

We need to pick: the implementation language, the agent framework, the
MCP SDK, the LLM provider abstraction, the database, the web stack, the
HTTP framework, the auth model, the deployment target, and the open
data sources we will integrate with.

## Decision

| Concern | Choice | License |
|---|---|---|
| Backend language | Go 1.23+ | BSD-3-Clause |
| Agent framework | `google.golang.org/adk` (Google ADK for Go) | Apache-2.0 |
| MCP SDK | `github.com/modelcontextprotocol/go-sdk` | Apache-2.0 |
| LLM abstraction | ADK Go's model layer + thin adapters | Apache-2.0 |
| Database | PostgreSQL 16 + PostGIS + pgvector + pg_trgm + uuid-ossp | PostgreSQL + Apache-2.0 + MIT-style |
| Query layer | `sqlc` (raw SQL → typed Go) | MIT |
| Migrations | `goose` | Apache-2.0 |
| HTTP router | `chi` | MIT |
| Templates | `templ` (`a-h/templ`) | MIT |
| Interactivity | HTMX + (optional) Alpine.js | BSD/MIT |
| Styling | Tailwind CSS (standalone CLI) | MIT |
| Maps | Leaflet 1.9.4 + OSM raster tiles | BSD-2-Clause |
| Observability | `slog` + OpenTelemetry Go + Prometheus client | Apache-2.0 / BSD |
| Auth | Signed-cookie sessions (web) + hashed API keys (API) | — |
| Deployment | Docker Compose (dev) + one Dockerfile per service (prod) | — |

## Rationale

### Why Go

- Google released **ADK for Go** with full multi-agent support and an
  Apache-2.0 licence. The Python and Java ADKs share the same design,
  so we can read each other's docs.
- The **official MCP Go SDK** is co-maintained with Google and is
  production-stable.
- Go's static typing, predictable concurrency, and single-binary deploy
  fit a system that fans out to many MCP servers in parallel.
- Author preference for scalable backends.

### Why ADK over hand-rolled orchestration

- Multi-agent primitives (`LLMAgent`, `SequentialAgent`, `ParallelAgent`,
  `LoopAgent`, `FunctionTool`) are already implemented and tested.
- MCP client integration is built in.
- Model-agnostic by design — Gemini, Anthropic, OpenAI, Ollama, etc.
  can be plugged in via a thin adapter without changing agent code.

### Why MCP

- Tool reuse: the same MCP server can be consumed by Joblantern's agent,
  by a CLI, by Claude Desktop, by Cursor, or by future federated
  Joblantern instances.
- Each MCP server is its own Go binary in `cmd/mcp-<name>/`, deployable
  and rate-limited independently.

### Why PostgreSQL + PostGIS + pgvector

- One engine for spatial, fuzzy text, and embeddings — fewer moving
  parts in an NGO deployment.
- PostGIS is the gold-standard open-source spatial database.
- `pgvector` covers our embedding scale comfortably with HNSW.
- `pg_trgm` handles fuzzy company-name and address matches without a
  separate search engine.

### Why server-rendered templ + HTMX

- The dominant production fullstack-Go pattern in 2025–2026.
- No npm build chain in the binary's critical path. Tailwind ships as
  a standalone binary; the app does not need Node.
- Accessibility, performance, and progressive enhancement come almost
  for free.

### Why we are explicitly **not** using

| Rejected | Reason |
|---|---|
| Google Maps / Mapbox SaaS for the critical path | Cost; commercial T&Cs; lock-in |
| Pinecone / proprietary vector DB | Cost; not needed at our scale |
| React / Next.js | Adds JS toolchain we do not need |
| Gin / Echo / Fiber | `chi` is idiomatic and composable |
| LiteLLM (Python) | Adds a Python runtime we do not need |
| Python in the runtime path | Single-binary deploy is a feature |
| wkhtmltopdf | LGPLv3 — not permissive enough for Joblantern |
| GPL / AGPL / LGPL / SSPL / BUSL / Elastic dependencies | Not Apache-2.0-compatible |

## Consequences

### Positive

- The whole stack is permissively licensed; downstream consumers can
  redistribute Joblantern under Apache 2.0 without surprises.
- Single Go binary per service simplifies deployment for NGOs.
- A reviewer can read the entire stack without context-switching
  languages.

### Negative

- Go's ML/embeddings ecosystem is thinner than Python's. We mitigate by
  using off-the-shelf embedding models exposed over HTTP (Ollama,
  hosted providers) rather than running model code in Go.
- `templ` requires a code-generation step. We accept this in exchange
  for type-safe templates.

### Open issues

- Multi-provider LLM routing complexity may eventually justify
  introducing **Bifrost** (Go-native, open-source) as a gateway. Not
  in v1.
- Long-form PDF export may need a richer renderer than `gofpdf`. We
  will revisit if real layout requirements emerge.

## References

- Apache License 2.0: <https://www.apache.org/licenses/LICENSE-2.0>
- Google ADK for Go: <https://github.com/google/adk-go>
- MCP Go SDK: <https://github.com/modelcontextprotocol/go-sdk>
- PostGIS: <https://postgis.net/>
- pgvector: <https://github.com/pgvector/pgvector>
- chi: <https://github.com/go-chi/chi>
- templ: <https://github.com/a-h/templ>
- HTMX: <https://htmx.org/>
- Leaflet: <https://leafletjs.com/>
