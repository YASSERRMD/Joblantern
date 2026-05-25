# Joblantern

[![CI](https://github.com/yasserrmd/joblantern/actions/workflows/ci.yml/badge.svg)](https://github.com/yasserrmd/joblantern/actions/workflows/ci.yml)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8.svg)](https://go.dev/dl/)
[![Status: pre-alpha](https://img.shields.io/badge/status-pre--alpha-orange.svg)](#status)

> **An agentic AI system that protects job seekers from fraudulent
> employers, fake listings, and exploitative recruitment schemes.**

Joblantern takes a job listing — a URL, a recruiter message, an offer
letter — and runs it through a multi-agent verification pipeline that
checks the company's registration, the domain's history, public scam
databases, salary plausibility, and recruitment-fee legality in the
destination country. Every signal cites its source and every verdict
is reproducible.

## Why

Migrant workers and first-time job seekers lose billions every year to
fraudulent recruiters, fake licensing fees, and trafficking-adjacent
schemes that look completely real on the surface. The tools to verify
a recruiter already exist and are mostly free and open. They are just
not stitched together. Joblantern stitches them together.

## Status

Pre-alpha. Active development.

## Stack

- **Backend:** Go 1.23+
- **Agent framework:** [Google ADK for Go](https://github.com/google/adk-go)
- **MCP:** [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- **Database:** PostgreSQL 16
- **Web:** chi + templ + HTMX + Tailwind CSS
- **Observability:** OpenTelemetry + Prometheus + slog

See [`docs/ADR/`](docs/ADR/) for locked architectural decisions.

## Quickstart

```bash
# clone
git clone https://github.com/yasserrmd/joblantern.git
cd joblantern

# environment
cp .env.example .env

# bring up postgres
make docker-up

# run migrations
make migrate-up

# build and run
make build
./bin/joblantern
```

Visit `http://localhost:8080/healthz` — you should see `ok`.

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
