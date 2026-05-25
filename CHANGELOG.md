# Changelog

All notable changes to Joblantern are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and the project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] — 2026-05-25

The first end-to-end release of Joblantern.

### Added

- **Scaffold (Phase 01):** Apache 2.0 licensing, Contributor Covenant 2.1 CoC,
  SECURITY policy, CI pipeline (fmt / vet / lint / test / build / license-check),
  Postgres compose, distroless runtime image, smoke test.
- **Database (Phase 02):** Custom Postgres image with PostGIS, pgvector,
  pg_trgm, citext. Goose migrations 0001–0010 covering users, sessions,
  verifications (with GIS geometry), evidence_facts, scam_reports (GIST +
  HNSW + GIN trgm), jurisdictions, mcp_audit_log. Sqlc-generated typed Go.
- **MCP foundation (Phase 03):** Shared `internal/mcpclient` wrapper with
  retry / timeout / audit hook. Reference `cmd/mcp-hello` server (stdio +
  streamable HTTP). `config/mcp.yaml` registry with loader + validator.
- **MCP servers (Phases 04–12):** Self-hosted Nominatim + Overpass deploy
  templates. Eight production MCP servers:
  - `joblantern.address` (Nominatim + Overpass classification)
  - `joblantern.streetview` (Mapillary)
  - `joblantern.registry` (OpenCorporates + provider interface)
  - `joblantern.domain` (WHOIS + crt.sh + Wayback)
  - `joblantern.scamdb` (Postgres-backed scam catalogue)
  - `joblantern.pattern` (YAML rule pack red-flag classifier)
  - `joblantern.salary` (bundled salary bands + FX)
  - `joblantern.law` (10-country recruitment-fee legality)
  - `joblantern.routing` (OpenRouteService)
- **Agent core (Phase 13):** Parallel sub-agent orchestrator,
  `/api/v1/verify` async API, in-memory store, two built-in sub-agents
  (pattern + language).
- **Risk engine (Phase 14):** Deterministic scoring with contradiction
  detection; agent uses it via `Orchestrator.WithScorer`.
- **Web UI (Phase 15):** Server-rendered HTML home form, async result page
  with auto-refresh, evidence list, risk badge, attribution footer.
- **i18n (Phase 16):** Embedded JSON catalogs for English, Arabic, Hindi,
  Tagalog, Bengali, Urdu, Swahili, Spanish; cookie + Accept-Language
  resolver; RTL direction metadata.
- **Feedback loop (Phase 17):** migration 0011 + `verification_feedback`
  queries, privacy + retention docs.
- **Observability (Phase 18):** Prometheus `/metrics` endpoint exposing
  `joblantern_http_requests_total` and `joblantern_http_request_duration_seconds`.
- **Hardening (Phase 19):** Security headers (CSP, HSTS, X-Frame-Options,
  X-Content-Type-Options, Referrer-Policy), per-IP rate limit on `/verify`,
  govulncheck in CI, data retention policy.

### Documentation

ADR-001 (stack), ADR-002 (schema), ADR-003 (MCP pattern), ADR-004
(self-host Nominatim), ADR-007 (no scraping), ADR-008 (rules vs LLM),
ADR-009 (law data cadence), ADR-010 (engine vs LLM for scoring), ADR-011
(feedback privacy). Per-server specs under `docs/MCP-SPECS/`.

### Not in scope (deferred to v0.2+)

- Live LLM provider integration (orchestrator and risk engine are
  provider-agnostic; LLM-driven narrative will land when at least one
  permissively-licensed provider adapter is wired).
- Mapbox-style map UI (PostGIS data is queryable; the UI does not yet
  render a Leaflet map).
- Templ-based HTML rendering with HTMX progressive enhancement.
- Container images published to ghcr.io.
- Browser extension and mobile PWA.
