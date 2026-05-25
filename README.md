# Joblantern

[![CI](https://github.com/yasserrmd/joblantern/actions/workflows/ci.yml/badge.svg)](https://github.com/yasserrmd/joblantern/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8.svg)](https://go.dev/dl/)
[![Status: pre-alpha](https://img.shields.io/badge/status-pre--alpha-orange.svg)](#status)

> **A GIS-driven agentic AI system that protects job seekers from fraudulent
> employers, fake listings, and exploitative recruitment schemes.**

Joblantern takes a job listing — a URL, a recruiter message, an offer letter —
and runs it through a multi-agent verification pipeline that interrogates the
claimed address, the company's registration, the domain's history, public scam
databases, salary plausibility, recruitment-fee legality in the destination
jurisdiction, and commute realism. Every signal cites its source, every
verdict is reproducible, and every map shows you what the agent actually saw.

## Why

Migrant workers and first-time job seekers lose billions every year to
fraudulent recruiters, fake licensing fees, and trafficking-adjacent
schemes that look completely real on the surface. The tools to verify a
recruiter — geocoders, registries, certificate transparency logs,
street-level imagery — already exist and are mostly free and open. They
are just not stitched together. Joblantern stitches them together.

## Status

Pre-alpha. Active development. Tracking the phased roadmap in
[`docs/CLAUDE_CODE_PROMPT_PACK.md`](docs/CLAUDE_CODE_PROMPT_PACK.md).

## Stack

- **Backend:** Go 1.23+
- **Agent framework:** [Google ADK for Go](https://github.com/google/adk-go)
- **MCP:** [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- **Database:** PostgreSQL 16 with PostGIS, pgvector, pg_trgm
- **Web:** chi + templ + HTMX + Tailwind CSS + Leaflet
- **Observability:** OpenTelemetry + Prometheus + slog

See [`docs/ADR/`](docs/ADR/) for locked architectural decisions.

## Quickstart

```bash
# clone
git clone https://github.com/yasserrmd/joblantern.git
cd joblantern

# environment
cp .env.example .env
# edit .env and fill in any required tokens

# bring up postgres
make docker-up

# run migrations (Phase 02+)
make migrate-up

# build and run
make build
./bin/joblantern
```

Visit `http://localhost:8080/healthz` — you should see `ok`.

## Repository layout

See [`docs/CLAUDE_CODE_PROMPT_PACK.md`](docs/CLAUDE_CODE_PROMPT_PACK.md)
for the full layout and phase plan.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

Every dependency must be Apache 2.0, BSD, MIT, or another permissive
license — no GPL, AGPL, LGPL, SSPL, BUSL, or Elastic License code may
ship in the Joblantern binary.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). All contributions are accepted
under the Developer Certificate of Origin (DCO).

## Security

See [SECURITY.md](SECURITY.md) for vulnerability disclosure.

## Author

Mohamed Yasser — Solutions Architect.
