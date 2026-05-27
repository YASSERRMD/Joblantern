# Research API — Changelog & Versioning Policy

The public research API follows a deliberate, dual-track versioning
policy so research output remains reproducible while the API evolves.

## Version identifiers

- REST: `/api/v{N}/...` (currently `v1`).
- GraphQL: a single `/graphql` endpoint that ships **named operation
  contracts**. Every breaking change increments a contract version
  exposed via the `__contract` meta field.

## Compatibility guarantees

| Change                                    | Allowed in same major | Notes |
|------------------------------------------|------------------------|-------|
| Adding a field to a response             | Yes                    | Clients must tolerate unknown fields. |
| Adding an enum value                     | Yes                    | Clients must default unknown values. |
| Renaming a field                         | No                     | Requires a new major. |
| Removing a field                         | No                     | Aliased + deprecated for ≥6 months first. |
| Tightening validation                    | No                     | Loosening is allowed. |

## Deprecation

A deprecated field is marked in OpenAPI with `deprecated: true` and in
GraphQL with `@deprecated`. The schedule is published in the public
[CHANGELOG](../CHANGELOG.md) entry for the release that introduced the
deprecation.

## Minimum support window

- **v1**: maintained until at least 24 months after `v2` ships.
- Every major receives **security patches** for 36 months after sunset
  announcement.

## Reproducibility

Annual archive snapshots (see [ARCHIVAL](ARCHIVAL.md), Phase 54)
include the API version that produced them and a frozen OpenAPI doc.

## Changelog format

Each release section in `CHANGELOG.md` has a `### API` subsection that
lists added, deprecated, and removed surface area.
