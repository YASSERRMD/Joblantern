# MCP Server — `joblantern.address`

## Purpose

Forward/reverse geocoding (via self-hosted Nominatim) and land-use
classification (via self-hosted Overpass). Consumed by the address
sub-agent to decide whether a claimed office address is real and
plausibly commercial.

## Identifier

- **Server name:** `joblantern.address`
- **Binary:** `cmd/mcp-address/`
- **Transport (dev):** `stdio`
- **Transport (prod):** `streamable HTTP` on `:8082`

## Tools

| Name | Purpose |
|---|---|
| `verify_address_exists` | Forward geocode an address; returns lat/lon + OSM ids. |
| `reverse_geocode` | Reverse geocode a coordinate. |
| `classify_building_type` | residential / commercial / mixed / unknown via Overpass tags. |
| `address_cluster_check` | (Stub in v1.) Count companies sharing this address in scam DB. |
| `_meta_attribution` | Required OSM attribution string for the UI. |

### Error codes

`ADDRESS_NOT_FOUND`, `RATE_LIMITED`, `UPSTREAM_TIMEOUT`, `INVALID_ARGS`.

## Rate limits

- Public Nominatim: max 1 req/sec — **dev only**. Set `NOMINATIM_URL` to
  a self-hosted instance for production.
- Public Overpass: best-effort; production must self-host
  (`deploy/overpass/`).

## License / attribution

- OpenStreetMap data is ODbL.
- Display "© OpenStreetMap contributors" wherever results are shown.
- Server exposes `_meta_attribution` for programmatic retrieval.
