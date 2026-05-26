# Transparency dashboard

`/transparency` exposes an anonymised public view of verdict
distribution across countries, suitable for journalists and policy
researchers.

## Privacy posture

- Every published row is `(date, country, risk) → count`. No
  identifying field, ever.
- Cells with raw count `< 5` are **dropped** entirely.
- Surviving counts are **fuzzed with discrete Laplace noise** (b ≈ 1)
  so a single-respondent contribution cannot be reverse-inferred from
  daily diffs. The `fuzzed` column flags this.
- Aggregates are computed on demand from the in-memory verification
  store in v0.1; a follow-up wires a nightly materialised view fed
  from `verifications` so the public endpoint stays cheap.

## Endpoints

| Path | Format |
|---|---|
| `GET /transparency` | HTML landing page |
| `GET /transparency/aggregate` | JSON array |
| `GET /transparency/aggregate.csv` | CSV |
| `GET /transparency/aggregate.jsonl` | newline-delimited JSON |

The HTML page intentionally avoids JS chart libraries: it's a card
total + table for everyone, no canvas / WebGL.

## Schema

```json
{
  "date": "2026-05-26",      // UTC, YYYY-MM-DD
  "country": "AE",           // ISO 3166-1 alpha-2, or "ZZ" if unknown
  "risk": "red",             // green | yellow | red
  "count": 12,               // potentially fuzzed
  "fuzzed": true
}
```

## Cache

The handler caches the rollup for 5 minutes. Long-running scrapes
should respect that cadence; a journalist-grade API tier with an opt-in
API key lands alongside the nightly job.

## ADR

Joblantern's overall data-minimisation stance is captured in
[`docs/PRIVACY.md`](PRIVACY.md). The transparency endpoint inherits
that stance and tightens it further via small-cell suppression and
DP noise.
