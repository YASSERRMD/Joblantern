# ADR-009 — Law data curation cadence

- **Status:** Accepted
- **Date:** 2025-11-09

## Decision

`data/recruitment-law.json` is the canonical Joblantern rule pack.
Every row carries a `citation_url` to an authoritative regulator
source. The pack is reviewed every 6 months. New countries are added
by amending the JSON and updating tests + spec.
