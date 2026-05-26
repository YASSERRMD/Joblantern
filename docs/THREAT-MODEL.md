# Joblantern threat model

Last reviewed: 2026-05-26. Review cadence: every 6 months and after
every shipped security-relevant phase.

## Assets

| Asset | Sensitivity |
|---|---|
| User-submitted listing text | Medium — may reveal what jobs the user is considering |
| Verification verdicts + evidence | Low — derived, but linkable to a user |
| User account email | High — directly identifying |
| Session tokens | High — grant account access |
| Recruiter API keys | High — grant the ability to issue badges |
| Federation private signing key | Critical — leaks let an attacker forge our peer signature |
| Recruiter-badge private signing key | Critical — same |
| Scam DB | Public-friendly, but small-corpus poisoning is a risk |

## Adversaries

### A1 — Network observer

Reads (and possibly modifies) traffic between the user and the
Joblantern instance.

- *Mitigations:* TLS everywhere, HSTS, security headers (Phase 19).
  Optional Tor hidden service (Phase 31). Lite mode minimises
  observable byte counts.
- *Residual:* TLS metadata still leaks "user contacted Joblantern".

### A2 — Hostile job seeker / scam-DB pollutor

Attempts to add false entries to the scam DB or game the rule pack.

- *Mitigations:* Direct scam-DB writes restricted to moderation
  (Phase 17). Pattern rules versioned and reviewed (ADR-008). No
  user-supplied YAML hot-loads.
- *Residual:* Operator must moderate the feedback queue.

### A3 — Compromised recruiter API key

Attacker mints badges they were not entitled to.

- *Mitigations:* `trust_level` per recruiter org (Phase 28); operator
  can flip a row to `suspended` to invalidate all future verifies.
  Badge revocation endpoint (in-memory v0.1).
- *Residual:* Already-issued unrevoked badges remain valid until
  expiry. Operator must explicitly revoke individual badge ids.

### A4 — Compromised federation peer

A peer signs and pushes a forged scam signal targeting an innocent
company.

- *Mitigations:* Per-peer trust labels (`vouched` / `observed` /
  `untrusted`) — operator can demote a peer instantly. Dedupe by
  fingerprint so one bad peer can't stack the verdict by itself.
- *Residual:* The labelled-bad signal still appears in `/a2a/recent`
  until aged out of the ring buffer.

### A5 — Local-host compromise

An attacker gains shell on the same machine the Joblantern binary runs
on.

- *Mitigations:* Distroless nonroot images. No on-disk session
  secrets in v0.1 (in-memory). Postgres password is sometimes the
  weakest link — operators MUST rotate the default in production.
- *Residual:* The federation and badge private keys, when loaded
  from disk in a future PR, must be 0600 and in a path outside the
  app's CWD.

### A6 — State-level adversary

Compels the operator to hand over data or to backdoor the build.

- *Mitigations:* Data minimisation (`docs/PRIVACY.md`,
  `docs/RETENTION.md`). Lite + panic-wipe routes give the user
  device-side hygiene. Reproducible builds let downstream consumers
  detect a backdoored binary by hash mismatch.
- *Residual:* Reproducible-build verification depends on at least
  one independent rebuilder; operators should publish their build
  manifests.

## Out of scope (today)

- Side-channel attacks on the cryptographic primitives. We rely on
  Go's `crypto/ed25519` and `crypto/tls` being constant-time.
- Supply-chain compromise of upstream Go modules. Partially mitigated
  by `govulncheck` in CI; complete defence requires a vendor + audit
  workflow we have not committed to.
- Physical-device seizure of a logged-in user's machine. Joblantern
  cannot do better than the browser there; we provide `/panic-wipe`
  as the best available knob.
