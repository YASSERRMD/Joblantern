# ADR-004 — Self-host Nominatim in production

- **Status:** Accepted
- **Date:** 2025-11-09

## Context

`mcp-address` geocodes addresses via Nominatim. OpenStreetMap
operates a public Nominatim endpoint, but its acceptable-use policy
forbids heavy use or resale. Joblantern must not depend on it in
production.

## Decision

- **Dev:** allowed to point at the public endpoint.
- **Prod:** MUST self-host using `deploy/nominatim/`. CI / staging must
  pin `NOMINATIM_URL` to a non-public instance.
- The default region (Monaco) keeps cold-start trivial for
  contributors. Production overrides via `PBF_URL`.

## Consequences

- Operators choose a region matching their user base.
- We ship a sized table so they can estimate cost up front.
- The address MCP returns `RATE_LIMITED` from the Nominatim policy
  guard if it ever sees HTTP 429.
