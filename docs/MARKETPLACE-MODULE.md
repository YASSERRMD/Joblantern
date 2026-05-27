# Marketplace Module

Same families of scammers operate across job offers, rental
listings, used-car listings, freelance gigs, and online marketplaces.
Phase 55 extends Joblantern's coverage so a single network running
all of them is detected as one cluster.

## Surfaces

- Schema: [marketplace_verifications](../migrations/0019_marketplace_verifications.sql)
- MCP server: [mcp-marketplace](../cmd/mcp-marketplace/main.go)
- Rule pack: [redflags](../internal/marketplace/redflags/rules.go)
- Cross-link: [crosslink](../internal/marketplace/crosslink/link.go)
- UI: [ui](../internal/marketplace/ui/flow.go)
- Integration with Phase 42 + 43: [integration](../internal/marketplace/integration/external.go)
- Consumer-protection partnerships: [partnerships](../internal/marketplace/partnerships/consumer.go)
- Operator API: [operator](../internal/marketplace/operator/api.go)

## Rule highlights

| Code                      | Severity | Meaning |
|---------------------------|----------|---------|
| `non-escrow-payment`      | 5        | Off-platform escrow (WU, gift card, crypto-only). |
| `advance-fee`             | 5        | Non-refundable upfront fee. |
| `shipping-via-seller`     | 4        | Shipping fee paid via the seller, not the platform. |
| `price-way-below-market`  | 3        | Asking price < 30% of median. |
| `off-platform-pressure`   | 2        | Seller refuses platform communication. |

## Legal

See [MARKETPLACE-LEGAL](MARKETPLACE-LEGAL.md) for per-platform and
per-jurisdiction posture.

## See also

- [RENTAL-MODULE](RENTAL-MODULE.md)
- [GRAPH-ANALYSIS](GRAPH-ANALYSIS.md)
- [RECRUITER-API](RECRUITER-API.md)
