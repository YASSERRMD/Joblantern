# ADR-002 — Database schema for Joblantern v1

- **Status:** Accepted
- **Date:** 2025-11-09
- **Deciders:** Mohamed Yasser (author)
- **Supersedes:** —
- **Superseded by:** —
- **Related:** ADR-001 (stack choices)

## Context

The Joblantern agent needs persistent storage for (a) the verification
requests it processes, (b) every evidence fact gathered, (c) the
final risk verdict, (d) a curated local catalogue of scam reports
used by the scam-db MCP server, (e) per-jurisdiction recruitment-law
rules used by the law MCP server, and (f) an append-only audit log of
every MCP tool call so verdicts are reproducible.

These workloads are simultaneously **spatial** (commute realism,
nearby-scam-cluster), **fuzzy text** (company-name and address
matching), **vector similarity** (scam-report semantic search), and
**ordinary relational** (users, sessions, verifications). We need to
decide whether to use one engine or several.

## Decision

Use a **single PostgreSQL 16 database** with the following extensions:

| Extension   | Purpose |
|---|---|
| `postgis`   | Spatial types (`geography(Point, 4326)`), distance, contains, intersects, GiST indexes |
| `vector`    | pgvector — embedding storage and HNSW cosine-similarity index |
| `pg_trgm`   | Trigram fuzzy matching (company names, addresses) |
| `uuid-ossp` | `uuid_generate_v4()` for primary keys |
| `citext`    | Case-insensitive text for emails |

The schema is owned by **goose-managed SQL migrations** in `migrations/`,
and the Go bindings are owned by **sqlc-generated code** in
`internal/db/`. Queries live in `queries/` as raw SQL.

## Tables (v1)

```
users (id, email, password_hash, created_at)

sessions (id, user_id → users, token_hash, expires_at, created_at)

verifications
  (id, user_id → users NULLABLE,
   raw_input jsonb, listing_url, recruiter_email, recruiter_phone,
   company_name, claimed_address, claimed_salary, role, jurisdiction,
   status enum-via-check, overall_risk enum-via-check, confidence,
   created_at, completed_at,
   claimed_geom geography(Point,4326))    [+ GiST index]

evidence_facts
  (id, verification_id → verifications,
   source, tool_name, fact_type controlled-vocab,
   value jsonb, supports_risk green/yellow/red/neutral,
   weight numeric in [0,1], citation, fetched_at)

scam_reports
  (id, company_name, address, address_geom geography(Point,4326),
   phone, email, domain,
   report_source, report_url, summary,
   embedding vector(384),
   reported_at, created_at)
  Indexes:
    GiST       (address_geom)
    GIN trgm   (company_name)
    HNSW       (embedding)
    btree      (phone), (email), (domain)

jurisdictions
  (code PK ISO 3166-1, name, recruitment_fee_legal,
   max_fee_pct, citation_url, updated_at)

mcp_audit_log
  (id, verification_id → verifications NULL,
   server, tool, args_hash sha256,
   latency_ms, status, error, called_at)
```

## Rationale

### Why one database, not many

NGOs and individual contributors will deploy Joblantern themselves. A
single engine that already handles spatial, fuzzy text, and embeddings
is dramatically easier to back up, migrate, and reason about than a
polyglot stack of Postgres + Elasticsearch + Pinecone + Redis. PostGIS
+ pgvector + pg_trgm meet every v1 requirement.

### Why `geography(Point, 4326)`

WGS84 lon/lat is universal and PostGIS's `geography` type gives metric
distances natively via `ST_DWithin(... , metres)` — no projection
gymnastics needed.

### Why HNSW over IVFFlat for embeddings

HNSW (pgvector ≥ 0.5) gives better recall for v1-scale corpora without
the IVFFlat training step. We use cosine distance because our embedding
model is unit-normalised.

### Why `pg_trgm` for company names

Trigram similarity tolerates typos, transliteration variants, and
abbreviations ("Acme Recruitment" vs "ACME Recruitments LLC"). It costs
a GIN index and is wildly faster than `LIKE '%…%'` scans.

### Why CHECK constraints, not native ENUM types

CHECK constraints can be evolved with simple `ALTER TABLE ... DROP/ADD
CONSTRAINT` — adding a new value to a native ENUM requires `ALTER TYPE
… ADD VALUE`, which interacts poorly with transactional migrations and
forces awkward gymnastics on rollback. Our vocabularies (`status`,
`overall_risk`, `supports_risk`, `fact_type`) all use CHECKs.

### Why JSONB for `raw_input`, `value`, etc.

The shape of a recruiter message or an MCP tool result will evolve over
time. JSONB lets us evolve without a schema migration; structured
fields can be promoted out of JSONB later if a query pattern emerges.

### Why an append-only audit log

A reviewer must be able to replay any verdict. `mcp_audit_log` records
every MCP call (server + tool + args hash + latency + status). It is
never updated, only inserted. Old rows can be rolled off into cold
storage on a TTL.

### Why `goose` + `sqlc` rather than an ORM

- Migrations are first-class, file-per-change, and reversible.
- Queries are owned as raw SQL — review-able, copy-pastable into psql,
  no ORM magic obscuring what hits the database.
- sqlc generates type-safe Go from those queries with zero runtime
  reflection.

## Consequences

### Positive

- Single Postgres engine. Single backup story. Single operational
  surface for NGO deployments.
- Schema and queries are reviewable artefacts (SQL files), not bytecode
  emitted by an ORM.
- PostGIS + pgvector are mature, permissively licensed, and run
  comfortably on commodity hardware.

### Negative

- pgvector's HNSW will tax memory at very large scale (millions of
  embeddings). Acceptable for v1; we can shard or move to a dedicated
  vector engine if and when scale demands it.
- sqlc's PostGIS support is limited (`geography` columns land as
  `interface{}`). We side-step this by always writing geography via
  `ST_SetSRID(ST_MakePoint(...), 4326)::geography` in the query and
  reading derived scalars (distance, bool) rather than the raw type.

### Open issues

- We may eventually want a `verification_feedback` table (Phase 17)
  to feed confirmed-scam outcomes back into the scam DB. That ships
  in a later migration when the feedback loop lands.

## References

- PostGIS: <https://postgis.net/>
- pgvector: <https://github.com/pgvector/pgvector>
- pg_trgm: <https://www.postgresql.org/docs/current/pgtrgm.html>
- goose: <https://github.com/pressly/goose>
- sqlc: <https://docs.sqlc.dev/>
