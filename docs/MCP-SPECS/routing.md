# MCP Server — `joblantern.routing`

## Purpose

Reachability / commute-realism checks via OpenRouteService.

## Tools

| Name | Purpose |
|---|---|
| `route` | Distance + duration between two coordinates. |
| `commute_realism_check` | `reachable` + `plausible` (duration < 3 h). |

### Error codes

`RATE_LIMITED`, `OUT_OF_REGION`, `INVALID_COORD`, `UPSTREAM_ERROR`.

## License / attribution

OpenRouteService free tier with API key (`ORS_API_KEY`). Attribute
routing results to OpenRouteService where surfaced in the UI.
