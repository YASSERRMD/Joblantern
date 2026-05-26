# Changelog

All notable changes to Joblantern are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and the project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] — 2026-05-26

Reach, deployment, channel, stakeholder, and community extensions on
top of v0.1.0. All work driven by `docs/ROADMAP-EXTENDED.md`.

### Added

- **Phase 21 — Browser extension.** Manifest v3 (Chrome) + v2 (Firefox)
  WebExtension under `ext/` with inline verdict badges on LinkedIn,
  Indeed, Bayt, Naukrigulf, GulfTalent, Jobstreet. Background service
  worker with IndexedDB cache. Popup + options. 8-language i18n
  parity. Build scripts produce signed-ready zips.
- **Phase 22 — PWA.** `manifest.webmanifest`, service worker with shell
  precache, offline fallback, background sync queue, web push handler.
  `/sw.js` and `/manifest.webmanifest` aliased to the application
  root.
- **Phase 23 — Mobile (Flutter).** `mobile/` Flutter project with home
  + result screens and a small HTTP client over the Joblantern API.
  Widget test verifies the home renders.
- **Phase 24 — Federation.** Migration `0012_federated_peers`. ed25519
  signed scam-signal envelopes. `/.well-known/joblantern.json` manifest,
  `POST /a2a/signal` ingest, `GET /a2a/recent`. Content-fingerprint
  dedupe. ADR-013.
- **Phase 25 — On-device inference.** `internal/llm` provider interface
  with Ollama implementation. Edge bundle Compose
  (`deploy/edge/docker-compose.yml`) for Raspberry Pi / mini-PC
  offline deployments.
- **Phase 26 — Telegram bot.** `cmd/joblantern-bot` with a transport-
  agnostic `botcore` (HTTP client, per-chat sessions, rate limit,
  formatting). Compose `bots` profile keeps it opt-in.
- **Phase 27 — More MCP servers.** `joblantern.registry.uk` (UK
  Companies House, OGL v3) and `joblantern.vies` (EU VAT validation).
  Demonstrates the established Phase 04 server template; remaining
  servers in the Phase 27 list are independently shippable.
- **Phase 28 — Recruiter signed badges.** Migration
  `0014_recruiter_badges`. ed25519 JWT-style tokens. `POST /api/v1/recruiter/badges`,
  `GET /badge/{token}`, `GET /badge/{token}/svg`.
- **Phase 29 — Transparency dashboard.** Anonymised aggregates with
  discrete-Laplace noise + small-cell suppression. `/transparency` +
  JSON / CSV / JSONL exports.
- **Phase 30 — NGO operator tools.** `/kiosk` large-font walk-in form +
  single-page A4 PDF printout at
  `/verifications/{id}/print.pdf` via gofpdf.
- **Phase 31 — Hostile-network hardening.** Tor onion service compose
  (`deploy/tor/`), `/lite` minimum-bandwidth form, `/panic-wipe` with
  Clear-Site-Data + IndexedDB / SW / cache wipe. New
  `docs/THREAT-MODEL.md` enumerating A1–A6 adversaries.
- **Phase 32 — Plugin architecture.** `joblantern-mcp.yaml` manifest
  schema. ed25519-signed manifests. Permissive-license allowlist
  (Apache, MIT, BSD, ISC, MPL). Migration `0013_plugins`.
- **Phase 33 — Voice interface.** `internal/voice` ASR (whisper.cpp HTTP)
  and TTS (piper) clients. Native binaries stay out-of-process so the
  single-binary distroless model is preserved.
- **Phase 34 — Continuous learning pipeline.** Anonymised CSV export +
  per-rule effectiveness scoring with `keep` / `review` /
  `consider_removal` recommendations. Human-in-the-loop only — never
  auto-applies changes.
- **Phase 35 — Sustainability docs.** `docs/GOVERNANCE.md`,
  `docs/STABILITY.md`, `docs/DEPRECATION.md`, `docs/FUNDING.md`.

### Deferred (intentionally) for v0.3

- Phase 23 follow-ups: drift offline cache, biometric lock, deep
  links, F-Droid metadata, CI APK/IPA builds (need iOS / Android
  signing infra).
- Phase 27 expansion: SEC EDGAR, POEA, MoHRE, BBB, Glassdoor,
  phone-reputation, ACLED, OpenSky.
- Phase 28 follow-ups: KYB onboarding flow, pricing tier, member
  roles, persistent badge signing key, admin revocation UI.
- Phase 30 follow-ups: IMAP listener, weekly digest email, white-
  label theming, case-management tagging, training mode, hotline
  call logger.
- Phase 32 follow-ups: plugin admin UI, sigstore signing, sandboxing
  ADR.
- Phase 33 follow-ups: `/api/v1/voice/verify` endpoint, Telegram
  voice-message handler.
- Phase 34 follow-ups: nightly aggregation cron, bias audit, drift
  detection, counterfactual evaluator, model cards.
- Phase 35 follow-ups (human-only): trademark policy, third-party
  security audit engagement, governance ratification by a public PR
  vote, Open Collective / GitHub Sponsors wiring, reference public
  instance.

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
