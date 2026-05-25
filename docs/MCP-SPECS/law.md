# MCP Server — `joblantern.law`

## Purpose

Per-jurisdiction recruitment-fee legality, licensing requirements, and
official regulator citation. Bundled JSON; mirrors
`data/recruitment-law.json`. Curated in v1 (see ADR-009).

## Tools

| Name | Purpose |
|---|---|
| `fee_legality_check` | Is the claimed fee legal in the country? |
| `check_recruiter_license` | v1 returns NOT_IMPLEMENTED + regulator URL. |
| `lookup_visa_requirements` | Returns regulator citation URL. |
| `_meta_disclaimer` | Mandatory legal disclaimer for the UI. |

### Error codes

`JURISDICTION_UNKNOWN`, `CITATION_MISSING`, `NOT_IMPLEMENTED`.

## Coverage (v1)

UAE, Saudi Arabia, Qatar, Bangladesh, Philippines, India, Pakistan,
Nepal, Sri Lanka, Indonesia.
