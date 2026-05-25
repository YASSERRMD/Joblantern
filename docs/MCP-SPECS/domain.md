# MCP Server — `joblantern.domain`

## Purpose

Compose port-43 WHOIS, crt.sh certificate transparency, and the
Internet Archive Wayback Machine into a single domain profile. Used
by the registry/domain sub-agent to spot brand-new domains that claim
long operation.

## Tools

| Name | Purpose |
|---|---|
| `domain_age` | Creation date, age in days, registrar, country (WHOIS). |
| `ssl_history` | Cert count, first/last issued, unique issuers (crt.sh). |
| `archive_history` | Earliest/latest snapshot + count (Wayback CDX). |
| `full_domain_profile` | Composite + derived `freshness_score` ∈ [0,1]. |

### Error codes

`INVALID_DOMAIN`, `WHOIS_UNAVAILABLE`, `RATE_LIMITED`, `UPSTREAM_ERROR`.

## License / attribution

- crt.sh — public data, free use.
- Internet Archive Wayback — attribution requested; display
  "Archive snapshots: Internet Archive" alongside archive results.
