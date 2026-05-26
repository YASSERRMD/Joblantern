# Joblantern — Extended Roadmap (Phases 21 → 35)

> Forward-looking plan covering reach, deployment, channel, data, stakeholder,
> security, community, accessibility, quality and sustainability extensions.
>
> See [`ROADMAP.md`](ROADMAP.md) for the short-form v0.2 / v0.3 picture and
> [`../CHANGELOG.md`](../CHANGELOG.md) for what shipped in v0.1.0.

Every phase below is **independently shippable**. They do not need to be done
in order, and many can run in parallel once Phase 32 (plugin architecture)
lands. Each phase follows the established Joblantern convention: a branch
named `phase/NN-name`, atomic commits, a PR per phase, merge to `main`,
optional version tag.

---

## Phase 21 — Browser extension (Chrome + Firefox)

**Branch:** `phase/21-browser-extension`
**Goal:** A WebExtension that adds a *"Verify with Joblantern"* button on
major job boards (LinkedIn, Indeed, Bayt, Naukrigulf, GulfTalent, Jobstreet)
and shows an inline verdict badge.

**Atomic commits:**

1. `ext`: scaffold manifest v3 webextension in `ext/` — shared codebase, browser-specific manifests.
2. `ext`: content script detects job listing on supported sites — DOM selectors per site.
3. `ext`: extract company name, address, salary, recruiter contact from page.
4. `ext`: popup UI with verdict summary.
5. `ext`: background script calls Joblantern `/api/v1/verify` — user-configurable endpoint and key.
6. `ext`: inline badge injected near job title — green/yellow/red dot with hover detail.
7. `ext`: settings page — API endpoint, language, opt-in telemetry.
8. `ext`: local cache of verdicts by listing URL — IndexedDB, 7-day TTL.
9. `ext`: Chrome Web Store + AMO packaging scripts — `make ext-package-chrome`, `make ext-package-firefox`.
10. `ext`: i18n parity with web app — share locale catalogs.
11. `ext`: privacy notice + permissions justification — required for store review.
12. `ext`: e2e tests with Playwright on a fixture HTML page per site.
13. `docs`: `docs/EXTENSION.md` — install, settings, supported sites.
14. `docs`: `ADR-012-extension-architecture.md`.

**Done:** Extension loads in Chrome dev mode, detects a listing on LinkedIn,
returns a verdict badge inline.

---

## Phase 22 — Progressive Web App (mobile)

**Branch:** `phase/22-pwa`
**Goal:** Offline-capable mobile experience using the existing server-rendered
UI plus a service worker.

**Atomic commits:**

1. `pwa`: add `manifest.webmanifest` — icons, splash, theme color.
2. `pwa`: add service worker with workbox-sw via CDN — precache static assets, runtime cache for verdicts.
3. `pwa`: offline fallback page.
4. `pwa`: add-to-home-screen prompt — smart timing.
5. `pwa`: background sync for queued submissions — submit when reconnected.
6. `pwa`: web push for verdict-ready notifications — VAPID keys via env.
7. `pwa`: push subscription persisted in db — migration `0017_push_subscriptions`.
8. `pwa`: server endpoint to send push on verdict completion.
9. `pwa`: lighthouse PWA score audit in CI — must stay ≥ 90.
10. `pwa`: install prompt copy localised in all 8 languages.
11. `pwa`: tests for offline flow with chromedp.
12. `docs`: `docs/PWA.md`.

**Done:** A user installs Joblantern to their home screen, submits a listing
offline, gets a push when the verdict arrives.

---

## Phase 23 — Native mobile (Flutter)

**Branch:** `phase/23-mobile-native`
**Goal:** A native mobile client for low-bandwidth regions where PWAs
underperform. Flutter chosen for shared iOS/Android codebase and offline-first
SQLite mirror.

**Atomic commits:**

1. `mobile`: Flutter project in `mobile/` — Dart 3, Riverpod state, Dio HTTP.
2. `mobile`: API client generated from OpenAPI spec — autogen.
3. `mobile`: home screen with submit flow.
4. `mobile`: progress screen via SSE.
5. `mobile`: verdict screen with map (`flutter_map` + OSM tiles).
6. `mobile`: offline-first verdict cache via `drift` (sqlite).
7. `mobile`: locale support matching web — Flutter `intl`.
8. `mobile`: dark mode + accessibility audit.
9. `mobile`: biometric lock for sensitive history — `local_auth`.
10. `mobile`: deep links from web verdict URLs — `uni_links`.
11. `mobile`: CI build for Android APK and iOS IPA.
12. `mobile`: F-Droid metadata — for reproducible builds.
13. `mobile`: tests in widget + integration.
14. `docs`: `docs/MOBILE.md`.

**Done:** Android APK installs, signs in, submits, displays verdict
offline-capable. iOS build green in CI.

---

## Phase 24 — Federated Joblantern instances (A2A)

**Branch:** `phase/24-federation`
**Goal:** Multiple NGO-operated Joblantern instances cooperate via ADK's A2A
protocol so a verdict from one instance enriches verdicts at others, without
sharing user PII.

**Atomic commits:**

1. `fed`: A2A server endpoint `/a2a` — exposes the orchestrator as a remote agent.
2. `fed`: A2A client wrapper — call remote Joblantern instances.
3. `fed`: instance registry table — migration `0018_federated_peers(url, name, trust_level, last_seen)`.
4. `fed`: signed instance manifests — ed25519 key per instance, manifest at `/.well-known/joblantern.json`.
5. `fed`: trust policy YAML — which peers count toward evidence, weight per peer.
6. `fed`: anonymised scam-signal exchange protocol — share company name + country + pattern tags, never user PII.
7. `fed`: dedup logic when same scam is reported by multiple peers.
8. `fed`: rate limits and abuse controls on incoming A2A calls.
9. `fed`: UI surface *"verified by N peer instances"*.
10. `fed`: tests with a 3-instance docker-compose.
11. `docs`: `docs/FEDERATION.md` — how an NGO joins the federation.
12. `docs`: `ADR-013-federation-trust-model.md`.

**Done:** Two instances running in compose exchange a scam signal; the
second instance's verdict includes the first instance's evidence.

---

## Phase 25 — On-device inference

**Branch:** `phase/25-on-device-llm`
**Goal:** For users in low-connectivity regions or with privacy concerns, run
a small local LLM and embedding model on-device or on an edge box.

**Atomic commits:**

1. `edge`: add `llama.cpp` runner support via local HTTP endpoint.
2. `edge`: Ollama provider adapter.
3. `edge`: bundled small LLM recommendation — Qwen2.5-3B or Phi-3-mini, all-MiniLM-L6-v2 for embeddings.
4. `edge`: quantisation documentation — Q4_K_M, Q5_K_M trade-offs.
5. `edge`: offline mode flag — disables non-essential external MCP calls.
6. `edge`: edge bundle Dockerfile — minimal image with Postgres + Ollama + Joblantern, runs on a 4 GB Raspberry Pi.
7. `edge`: model download wizard on first run.
8. `edge`: degraded-evidence labels in UI.
9. `edge`: power-aware scheduling — skip heavy MCP calls on battery.
10. `edge`: tests with mock local LLM.
11. `docs`: `docs/EDGE-DEPLOYMENT.md` — Raspberry Pi recipe.
12. `docs`: `ADR-014-on-device-inference.md`.

**Done:** A Raspberry Pi 5 with 8 GB RAM runs the full Joblantern stack offline
using Qwen2.5-3B and returns verdicts in under 60 s.

---

## Phase 26 — Telegram & WhatsApp bots

**Branch:** `phase/26-messaging-bots`
**Goal:** Many target users live inside Telegram and WhatsApp. Meet them there.

**Atomic commits:**

1. `bot`: `cmd/joblantern-bot` scaffolding — single binary, multiple adapters.
2. `bot`: Telegram adapter using `go-telegram-bot-api`.
3. `bot`: WhatsApp adapter via `whatsmeow`.
4. `bot`: conversation state machine — collect listing → submit → return verdict.
5. `bot`: rate limits per chat.
6. `bot`: locale auto-detection from Telegram user lang.
7. `bot`: image OCR for listings sent as screenshots — Tesseract via `gosseract`.
8. `bot`: privacy commands — `/forget`, `/export`.
9. `bot`: admin commands — `/stats`, `/health`.
10. `bot`: deploy as separate service in compose.
11. `bot`: tests with fake Telegram API.
12. `docs`: `docs/BOTS.md`.

**Done:** Sending a job-listing screenshot to the Telegram bot returns a
verdict reply with the verdict URL.

---

## Phase 27 — More MCP servers (expansion pack)

**Branch series:** `phase/27a-mcp-X`, `phase/27b-mcp-Y`, …
**Goal:** Broaden evidence sources without inflating any single phase.

Each new server follows the established MCP commit template (scaffold,
client, tools, cache, error codes, tests, OTel, Dockerfile, compose entry,
spec, ADR), ~14 commits per branch.

Suggested servers:

- `mcp-companies-house-uk` — UK Companies House (OGL).
- `mcp-sec-edgar-us` — US SEC EDGAR for publicly traded employers.
- `mcp-eu-vies` — EU VAT number validation.
- `mcp-poea-philippines` — POEA licensed-recruiter check.
- `mcp-mohre-uae` — UAE recruiter licensing.
- `mcp-osint-people` — minimal, ethics-bound: public press + Wikidata on claimed executives.
- `mcp-phone-reputation` — open phone-reputation feeds (PhoneInfoga-style).
- `mcp-bbb-scam-tracker` — if their open feed remains accessible.
- `mcp-glassdoor-public-reviews` — only public, scrape-allowed surfaces.
- `mcp-air-quality` — for *"hazardous-air industrial workplace not disclosed"* checks.
- `mcp-protest-risk` — ACLED conflict-event data.
- `mcp-flight-implausibility` — OpenSky historical to check whether *"you must fly tomorrow"* is feasible.

---

## Phase 28 — Recruiter-side API

**Branch:** `phase/28-recruiter-api`
**Goal:** Let legitimate recruiters and job boards pre-verify their own
listings before publishing, generating a public *"Verified by Joblantern"*
badge.

**Atomic commits:**

1. `recruiter`: db migration `0019_recruiter_orgs` — orgs, members, API keys.
2. `recruiter`: onboarding flow — KYB via OpenCorporates + manual review.
3. `recruiter`: API endpoint `POST /api/v1/recruiter/listings`.
4. `recruiter`: signed verdict badges — JWT badge embeddable on boards.
5. `recruiter`: public verifier endpoint `GET /badge/:id`.
6. `recruiter`: webhook callbacks on verdict change.
7. `recruiter`: monthly trust-score per recruiter org.
8. `recruiter`: pricing tier scaffolding (free for NGOs and low-income regions).
9. `recruiter`: ToS for recruiters.
10. `recruiter`: dashboard for recruiters to see listings + scores.
11. `recruiter`: tests for badge signing/verification.
12. `docs`: `docs/RECRUITER-API.md`.

**⚠ Caveat:** Building this turns Joblantern into a two-sided product. Worth
a deliberate decision before starting.

**Done:** A recruiter org signs up, submits a listing, receives a signed
badge, embeds it on their site, badge verifies publicly.

---

## Phase 29 — Public transparency dashboard

**Branch:** `phase/29-public-dashboard`
**Goal:** A public, anonymised dashboard showing aggregate scam trends per
country, top fraudulent company-name patterns, top abused jurisdictions,
year-over-year change. For journalists and policymakers.

**Atomic commits:**

1. `dash`: nightly aggregation job — materialised views.
2. `dash`: differential privacy noise on small cells.
3. `dash`: public landing page `/transparency`.
4. `dash`: country drill-down pages.
5. `dash`: time-series charts via server-side SVG.
6. `dash`: choropleth map of scam density.
7. `dash`: CSV + JSONL public download endpoints — versioned daily.
8. `dash`: API keys for journalists/researchers with higher rate limits.
9. `dash`: data dictionary docs.
10. `dash`: tests for aggregation correctness on a known fixture.
11. `dash`: RSS feed of monthly trend reports.
12. `docs`: `docs/TRANSPARENCY.md`.

**Done:** `/transparency` shows verdict distribution by country with safe
aggregation, a journalist can download the monthly CSV.

---

## Phase 30 — NGO operator tools

**Branch:** `phase/30-ngo-tools`
**Goal:** Tools for NGOs running Joblantern instances on behalf of vulnerable
populations.

**Atomic commits:**

1. `ngo`: case management mode — caseworkers tag and follow verdicts attached to specific individuals (with consent).
2. `ngo`: bulk import of suspicious listings from email forwards — IMAP listener.
3. `ngo`: weekly digest emails for caseworkers.
4. `ngo`: kiosk mode — large-text UI, no login.
5. `ngo`: paper-friendly verdict printout — single-page PDF in user's language.
6. `ngo`: hotline call logger — manual entry for users who call in.
7. `ngo`: trauma-informed UI variant — softer language, no alarming red flashes.
8. `ngo`: caseworker training mode — sample listings with explanations.
9. `ngo`: audit log export for accountability reporting.
10. `ngo`: white-label theming per NGO.
11. `ngo`: tests for case management flow.
12. `docs`: `docs/NGO-OPERATIONS.md`.

**Done:** An NGO caseworker can log in, see all cases for their clinic, print
a verdict in Tagalog for a walk-in client.

---

## Phase 31 — Hardening for high-risk regions

**Branch:** `phase/31-hostile-network`
**Goal:** Joblantern remains usable where networks are monitored, censored, or
unreliable.

**Atomic commits:**

1. `hostile`: domain-fronting deployment guide (when allowed).
2. `hostile`: Tor hidden service deployment recipe.
3. `hostile`: minimum-bandwidth mode — text-only verdicts, no map tiles.
4. `hostile`: encrypted-at-rest session data — additional layer beyond TLS.
5. `hostile`: panic-wipe button — clears local history instantly.
6. `hostile`: dummy-traffic mode.
7. `hostile`: alternative tile sources — cached common-region tile bundles.
8. `hostile`: IPFS distribution of static assets — experimental.
9. `hostile`: signed reproducible builds documentation.
10. `hostile`: Snowflake-style proxy support documentation.
11. `hostile`: threat model document.
12. `docs`: `docs/HOSTILE-NETWORK.md` + `docs/THREAT-MODEL.md`.

**Done:** A documented recipe exists for running Joblantern over Tor with
reproducible builds and panic-wipe.

---

## Phase 32 — Plugin architecture & MCP marketplace

**Branch:** `phase/32-plugin-architecture`
**Goal:** Let the community contribute MCP servers without forking Joblantern.

**Atomic commits:**

1. `plugin`: MCP server manifest spec — `joblantern-mcp.yaml`.
2. `plugin`: signed manifests — sigstore signing.
3. `plugin`: plugin registry — table + UI to add/remove third-party servers.
4. `plugin`: sandboxing recommendation — run untrusted servers in gVisor or Firecracker.
5. `plugin`: rating + review system for plugins.
6. `plugin`: official curation list — vetted, signed, recommended.
7. `plugin`: contribution guide for plugin authors.
8. `plugin`: example community plugin in `examples/community-mcp/`.
9. `plugin`: CI scanner for plugin manifests.
10. `plugin`: telemetry for plugin usage — opt-in.
11. `plugin`: tests for plugin registration + deregistration.
12. `docs`: `docs/PLUGINS.md`.

**Done:** A community member can publish a signed MCP server manifest, a
deployment can add it via the UI, and it shows up in agent decisions.

> **Strategic note:** Doing Phase 32 *before* Phase 27 lets third parties
> contribute MCP servers via the marketplace instead of forcing the core
> team to maintain ten more in `cmd/`.

---

## Phase 33 — Multilingual voice interface

**Branch:** `phase/33-voice`
**Goal:** Many migrant workers cannot read fluently in their second language
but can speak comfortably. Voice unlocks the product for them.

**Atomic commits:**

1. `voice`: `whisper.cpp` integration for transcription — local-first.
2. `voice`: TTS via `piper` or `coqui-tts` — local-first.
3. `voice`: voice submit endpoint — upload audio, transcribe, verify.
4. `voice`: voice verdict playback — TTS the verdict in user's language.
5. `voice`: Telegram bot voice message handler — extends Phase 26.
6. `voice`: IVR-style phone integration — optional, via Twilio or open-source PBX.
7. `voice`: dialect tuning notes — Emirati Arabic, Philippine English, Indian English, Bengali.
8. `voice`: tests with audio fixtures.
9. `voice`: privacy notes for voice data.
10. `voice`: bandwidth-aware codec choice.
11. `voice`: degraded transcription confidence shown in UI.
12. `docs`: `docs/VOICE.md`.

**⚠ Caveat:** Pulls in `llama.cpp`-class binary dependencies that break the
current single-Go-binary distroless deploy.

**Done:** A user records *"Is this Dubai job at XYZ Trading legit?"* in
Tagalog, gets a spoken verdict in Tagalog.

---

## Phase 34 — Continuous learning pipeline

**Branch:** `phase/34-learning-pipeline`
**Goal:** Improve rules, embeddings, and patterns from accumulated verdict
outcomes — without retraining a foundation model.

**Atomic commits:**

1. `ml`: labelled dataset builder from verdict + feedback — exports anonymised examples.
2. `ml`: rule effectiveness scoring — track which rules correlate with confirmed scams.
3. `ml`: weekly rule-pack proposal job — suggests rule weights based on outcomes.
4. `ml`: human review queue for proposed rule changes — never auto-applied.
5. `ml`: embedding model fine-tuning recipe — sentence-transformers domain adaptation.
6. `ml`: A/B testing harness for rule packs.
7. `ml`: drift detection — alerts when verdict distribution shifts unexpectedly.
8. `ml`: counterfactual evaluator — replay historical verdicts against new rule packs.
9. `ml`: model card generator for each release of the embedding model.
10. `ml`: bias audit pipeline — verdict distribution by country/language/etc.
11. `ml`: tests for the pipeline.
12. `docs`: `docs/LEARNING-PIPELINE.md` + `docs/MODEL-CARDS/`.

**Done:** Running the pipeline on a month of verdicts produces a rule-pack
proposal PR for human review.

---

## Phase 35 — v1.0.0 release & sustainability

**Branch:** `phase/35-v1-release`
**Goal:** Get Joblantern to a stable, maintained, sustainable v1.0.0.

**Atomic commits:**

1. `release`: v1 release notes — consolidated CHANGELOG.
2. `release`: stability guarantees document.
3. `release`: backwards-compatibility test suite.
4. `release`: deprecation policy.
5. `release`: governance document — maintainers, decision-making, contributor ladder.
6. `release`: funding documentation — Open Collective or GitHub Sponsors setup.
7. `release`: trademark policy for "Joblantern" name.
8. `release`: third-party security audit recommendations.
9. `release`: CVE disclosure process refinement.
10. `release`: production reference deployment — public instance for journalists.
11. `release`: anniversary blog post template.
12. `release`: cut `v1.0.0` tag — multi-arch images, signed release artifacts via sigstore.

**⚠ Human-only items.** Trademark policy, OSS audit engagement, funding setup,
and governance ratification are paperwork that needs human signatures and
should not be executed by an agent. The remaining commits (release notes,
test suite, CVE process, reference deployment, blog template, tag) are
appropriate to automate.

**Done:** v1.0.0 is tagged, signed, published. A reference instance is live.
A governance model is in place so the project can outlive any single
contributor.

---

## How these phases stack

| Phase range | Theme |
|---|---|
| 21–22 | Reach extensions: browser, mobile web. |
| 23–25 | Deployment extensions: native mobile, federation, edge. |
| 26 | Channel extension: messaging. |
| 27 | Data-source breadth. |
| 28–30 | Stakeholder extensions: recruiters, public, NGOs. |
| 31 | Threat-model extension. |
| 32 | Community extension. |
| 33 | Accessibility extension. |
| 34 | Quality extension. |
| 35 | Sustainability extension. |

## Open questions before execution

Before scheduling any of these phases, the following call-outs from v0.1.0
should be resolved:

1. **DB-backed `agent.Store`** to replace the in-memory store — currently
   every restart loses verdict history. Pick a phase to host this work
   (most natural: an early v0.2 follow-up before Phase 21).
2. **Real LLM provider wiring** — the Phase 13 `ScoreFunc` hook is in place
   but nothing calls a model yet. Pick a provider and adapter.
3. **Recruiter-side API (Phase 28)** is a product-shape decision — see the
   ⚠ caveat in that phase.
4. **Voice interface (Phase 33)** breaks the single-binary deploy. Decide
   whether the trade-off is acceptable before starting.
