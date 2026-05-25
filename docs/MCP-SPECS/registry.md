# MCP Server — `joblantern.registry`

## Purpose

Look up companies in a business registry and surface registration
status, age, registered address, and officers. Default provider:
OpenCorporates.

## Identifier

- **Server name:** `joblantern.registry`
- **Binary:** `cmd/mcp-registry/`
- **Transport (dev):** `stdio` | **Transport (prod):** HTTP `:8084`

## Tools

| Name | Purpose |
|---|---|
| `lookup_company` | Search by name (+ optional jurisdiction). |
| `get_company` | Fetch full record by provider id (e.g. `gb/12345`). |
| `check_registration_status` | Derive `is_active`, `is_recent`, `age_days`. |
| `_meta_attribution` | Returns "Powered by OpenCorporates". |

### Error codes

`COMPANY_NOT_FOUND`, `JURISDICTION_UNKNOWN`, `RATE_LIMITED`,
`TOKEN_INVALID`, `UPSTREAM_ERROR`, `INVALID_ARGS`.

## License / attribution

OpenCorporates data is ODbL with **share-alike** and a mandatory
"Powered by OpenCorporates" attribution. Display this attribution in
the footer of any UI that surfaces registry results.
