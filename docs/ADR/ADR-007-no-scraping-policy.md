# ADR-007 — No scraping of third-party scam catalogues

- **Status:** Accepted
- **Date:** 2025-11-09

## Decision

For v1, Joblantern's scam_reports table is populated only by:

1. Hand-curated entries committed to `data/seed-scams.csv`.
2. Moderated user feedback (Phase 17), which writes through
   `mcp-scam-db.insert_report`.

We **do not** scrape any upstream consumer-protection catalogue, BBB
listing, news article, or similar third-party source whose terms of
use prohibit derivative redistribution.

## Rationale

- Several attractive sources (e.g. BBB scam tracker, regulator
  blacklists) carry non-commercial or share-alike licences that would
  contaminate Joblantern's Apache-2.0 posture.
- Trust matters more than coverage for v1. A small, accurate corpus
  with citations beats a large corpus of dubious provenance.

## Consequences

- Coverage starts small (tens of seed rows).
- Future ADRs may approve specific sources whose terms permit
  redistribution.
