# MCP Server — `joblantern.scamdb`

## Purpose

Read/write surface for Joblantern's curated catalogue of recruitment
scams (`scam_reports` table). All scam-DB writes go through this
server so audit and moderation stay centralised.

## Tools

| Name | Purpose |
|---|---|
| `search_reports_by_phone` | Exact normalised-phone match. |
| `search_reports_by_email` | Exact email + `%@domain` form. |
| `search_reports_by_domain` | Exact domain + any subdomain. |
| `search_reports_by_company_name` | Trigram fuzzy match (`min_sim` default 0.4). |
| `search_reports_near_point` | PostGIS `ST_DWithin` around (lat, lon). |
| `insert_report` | Internal-only; moderation pipeline writes here. |

### Error codes

`INVALID_ARGS`, `DB_ERROR`, `EMBEDDING_UNAVAILABLE`.

## Ingestion

`scripts/seed-scam-db.sh` loads `data/seed-scams.csv`. New reports
should be added through the v1 manual-curation flow only — no
scraping of upstream catalogues whose terms forbid it (see ADR-007).
