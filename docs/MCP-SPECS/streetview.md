# MCP Server — `joblantern.streetview`

## Purpose

Pull recent Mapillary street-level imagery near a coordinate and report
when (or if) the location has visual ground-truth. Consumed by the
address sub-agent.

## Identifier

- **Server name:** `joblantern.streetview`
- **Binary:** `cmd/mcp-streetview/`
- **Transport (dev):** `stdio`
- **Transport (prod):** `streamable HTTP` on `:8083`

## Tools

| Name | Purpose |
|---|---|
| `images_near_point` | Up to 100 images within ≤50 m of (lat, lon). |
| `latest_image_age` | Age in days of the most recent image. `-1` = none. |
| `_meta_attribution` | Required Mapillary attribution. |

### Error codes

`TOKEN_INVALID`, `RATE_LIMITED`, `UPSTREAM_ERROR`, `INVALID_ARGS`.

## License / attribution

Mapillary imagery is CC BY-SA 4.0. Display
"Imagery © Mapillary (CC BY-SA 4.0)" alongside any thumbnail.
