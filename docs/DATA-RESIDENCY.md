# Data Residency

Tenant data lives in exactly one region and never crosses regions
even for operational convenience.

## Regions

| ID    | Ingress URL                  | Postgres primary |
|-------|------------------------------|-------------------|
| eu    | https://eu.joblantern.org    | Frankfurt |
| us    | https://us.joblantern.org    | Ashburn |
| apac  | https://ap.joblantern.org    | Singapore |

## Tenant residency choice

Tenants pick a region at signup. We default the choice by the
tenant's country of incorporation but allow overrides. Overrides
require a one-line justification stored on the tenant record.

## Cross-region exceptions

There are none for tenant data. Aggregated, anonymised research
exports can be published globally — see [RESEARCH-API](RESEARCH-API.md).

## Audit

Every cross-region call is rejected at the ingress edge. Attempts
are logged for compliance audit per the Phase 48 register.
