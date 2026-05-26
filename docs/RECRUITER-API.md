# Recruiter API (signed badges)

Legitimate recruiters and job boards can pre-verify their listings and
get a public, signed badge they can embed.

> Scope note: v0.1 is intentionally minimal — no pricing tiers, no
> member-management UI, no KYB onboarding flow. The trust label is set
> manually by the operator on the `recruiter_orgs` row (`vouched` /
> `observed` / `suspended`). See [`docs/ROADMAP-EXTENDED.md`](ROADMAP-EXTENDED.md)
> for the full Phase 28 scope.

## Endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/v1/recruiter/badges` | Issue a badge from a completed verification |
| `GET`  | `/badge/{token}`            | Public verifier — returns Claims as JSON |
| `GET`  | `/badge/{token}/svg`        | Renders a small SVG you can `<img src=…>` embed |

## Issue request

```json
{
  "org_id": "uuid-of-recruiter-org",
  "org_name": "Acme Recruiting",
  "verification_id": "uuid-of-completed-verification",
  "risk": "green",
  "trust_level": "vouched"
}
```

## Issue response

```json
{
  "token": "<opaque base64url token>",
  "claims": {
    "badge_id": "...", "org_id": "...", "org_name": "...",
    "verification_id": "...", "risk": "green",
    "issued_at": "2026-05-26T...", "expires_at": "2026-08-24T...",
    "issuer": "https://joblantern.example/",
    "trust_level": "vouched"
  }
}
```

## Verifier response

```json
{
  "valid": true,
  "claims": { ... }
}
```

Invalid (bad signature, expired, revoked) → HTTP 401 + `{"valid": false, "error": "..."}`.

## Embedding

```html
<a href="https://joblantern.example/badge/<token>">
  <img src="https://joblantern.example/badge/<token>/svg" alt="Verified by Joblantern">
</a>
```

## Signing key

v0.1 generates a fresh ed25519 keypair per `joblantern` process. **This
means existing badges fail to verify after a restart.** Production
must load a persistent key from disk (`JOBLANTERN_BADGE_KEY` env var
points at a file — follow-up PR wires the loader). The verifier always
uses the in-memory pubkey, so once persistent key loading lands no
client-side change is needed.

## Revocation

`BadgeIssuer.Revoke(badgeID)` flips an in-memory flag. A revoked
badge's verifier endpoint returns 401 even though the signature is
valid. The admin UI to revoke is a follow-up PR.
