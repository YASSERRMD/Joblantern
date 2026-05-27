# Research API

Joblantern exposes a public, anonymized research API so journalists,
academics, and migrant-rights NGOs can study recruitment-fraud
patterns without raw, identifying data ever leaving the trust
boundary.

## Surfaces

- **REST**: `/api/v1/*` — OpenAPI 3.1 spec served at `/api/v1/openapi.json`
- **GraphQL**: `/graphql` — read-only, with an explorer at `/graphql/explorer`
- **Bulk exports**: cursor-paginated, see [bulk.go](../internal/researchapi/export/bulk.go)
- **Webhooks**: HMAC-signed deliveries for new high-confidence verdicts

## Tiers

| Tier        | Auth needed                                      | Rate     | DUA |
|-------------|--------------------------------------------------|----------|-----|
| Public      | None                                             | 30 rpm   | No  |
| Academic    | Bearer token, institutional email verified       | 600 rpm  | Yes |
| Journalist  | Bearer token, editor counter-signature           | 300 rpm  | Yes |
| Regulator   | Phase 37 regulator integration                   | 1200 rpm | Yes |

## Data-use agreement (DUA)

The DUA flow is automated (no DocuSign dependency). Researchers sign
the agreement body with their own ed25519 key, Joblantern
counter-signs, and the resulting token embeds the agreement ID.

## Worked example

```bash
# 12 months of anonymized verdicts, GraphQL
curl -s https://joblantern.org/graphql \
  -H "authorization: Bearer $JL_TOKEN" \
  -H "content-type: application/json" \
  -d '{"query":"{ verdicts(since:\"2025-05-27\", first:1000) { id country industry riskBand } }"}'
```

## Anonymization commitments

- No raw company names, addresses, phone numbers, or applicant fields
- All free-text excerpts are model-paraphrased to remove identifiers
- Country and industry are k-anonymized: buckets < 5 are coalesced

## Limits

See [API-VERSIONING](API-VERSIONING.md) for deprecation policy, and
[ratelimit](../internal/researchapi/ratelimit/tier.go) for per-tier
caps. Query complexity is bounded — see
[graphql/cost.go](../internal/researchapi/graphql/cost.go).

## Stability

Once a field appears on this surface, it follows the deprecation
policy. New fields are additive within a major.
