# Federation

Multiple Joblantern instances can exchange anonymised scam signals. No
user-identifying data ever crosses an instance boundary.

## What's exchanged

A `Signal` carries:

- `origin_url` — the peer that produced it
- `company_name`, `country` — what the verdict was about
- `pattern_codes` — high-level pattern tags from the rule engine
- `verdict` — `"red"` or `"yellow"`
- `issued_at`

And nothing else. No raw listing text, no recruiter contact details, no
user identifiers.

## How peers trust each other

- Each instance generates an ed25519 keypair at first boot and
  publishes its public key at `/.well-known/joblantern.json`.
- Operators add peers via the `federated_peers` table (migration 0012)
  with a `trust_level`:
  - `vouched` — full weight in the local risk engine
  - `observed` — half weight; surface to user but don't dominate
  - `untrusted` — display-only
- Every incoming `POST /a2a/signal` envelope is verified against the
  stored peer pubkey. Unknown peers → 403. Bad signature → 401.

## Endpoints

| Path | Purpose |
|---|---|
| `GET /.well-known/joblantern.json` | This instance's signed manifest |
| `POST /a2a/signal` | Ingest a signed Envelope from a known peer |
| `GET /a2a/recent` | List recently-ingested signals (for the UI surface) |

## What's not yet done

- Outgoing push (a small worker that promotes confirmed-scam verdicts
  into Signals and POSTs them to enabled peers). The data types,
  signing, ingest, and dedupe are ready; the cron worker is a small
  follow-up PR.
- Persistence of ingested signals into Postgres. v1 keeps a 1024-entry
  in-memory ring. Migration to `scam_reports` cross-references is a
  natural next step.
- UI surface: *"verified by N peer instances"* badge on the result
  page.

## ADR

See [ADR-013](ADR/ADR-013-federation-trust-model.md).
