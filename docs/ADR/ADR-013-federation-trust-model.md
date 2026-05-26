# ADR-013 — Federation trust model

- **Status:** Accepted
- **Date:** 2026-05-26

## Decision

- **Pairwise trust**, not transitive. Each instance maintains its own
  `federated_peers` list. Adding peer X to your list does not import
  X's peers.
- **Three trust levels** (`vouched` / `observed` / `untrusted`), each
  mapped to an evidence weight by the local risk engine.
- **ed25519 signatures** on every signal. Verification uses the
  pubkey stored at registration time, not a freshly-fetched manifest
  (manifest fetch is for *registering* a peer; per-signal verification
  must not depend on the peer being online).

## Why pairwise + non-transitive

Joblantern instances are operated by NGOs in different jurisdictions
with very different legal exposures. An NGO must be able to vouch for
a peer it knows without auto-trusting whomever that peer trusts.

## Why ed25519

- Short keys (32 bytes), short signatures (64 bytes).
- Deterministic — same input always yields same signature, easier to
  test.
- In Go's standard library (`crypto/ed25519`), no extra dependency.

## Anonymisation

The Signal type is intentionally narrow. The risk engine produces much
more information than this; the federation layer deliberately discards
all of it except: company name, country, high-level pattern tags,
verdict, origin URL, issued-at timestamp.

## Consequences

- Operators can join the federation without exposing their users.
- A misbehaving peer's signals can be downgraded to `untrusted`
  without breaking the verification flow.
- Dedupe is by content-fingerprint, so a popular scam reported by
  five peers does not stack into a single overweighted verdict.
