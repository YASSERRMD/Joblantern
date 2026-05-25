# MCP Server — `joblantern.salary`

## Purpose

Compare a claimed salary to public salary bands for the role and
country. Ships with a bundled JSON dataset (`data/salary-bands.json`)
that mirrors the embedded `cmd/mcp-salary/bands.json`.

## Tools

| Name | Purpose |
|---|---|
| `salary_range_check` | Returns `within_range`, percentile, multiplier vs p50. |
| `currency_normalize` | Converts between bundled currencies. |

### Error codes

`UNKNOWN_ROLE`, `UNKNOWN_COUNTRY`, `CURRENCY_UNAVAILABLE`.

## Data sources

v1 ships a small curated set sourced from public ILO/ILOSTAT and
national statistics offices. Update cadence: every 6 months.
