# MCP Server — `joblantern.vies`

EU VAT-number validation against the official VIES service.

## Tool

| Name | Purpose |
|---|---|
| `validate_vat_number` | `{country, vat_number}` → `{valid, name, address}` |

## Error codes

`INVALID_ARGS`, `INVALID_COUNTRY`, `VAT_NOT_FOUND`, `UPSTREAM_ERROR`.

## Attribution

VIES is a free service from the European Commission. No attribution
contract; courtesy mention "Verified via EU VIES" alongside results.
