# Stability guarantees

What changes safely between releases, and what does not.

## API-stable surfaces (semver-respected)

Breaking changes require a major version bump (`v1.x.y` → `v2.0.0`).

- `POST /api/v1/verify` request shape (the `Submission` JSON).
- `GET /api/v1/verifications/{id}` response shape (`Record` JSON).
- `POST /api/v1/recruiter/badges` request + `GET /badge/{token}` response.
- `POST /a2a/signal` envelope wire format and signature algorithm.
- `joblantern-mcp.yaml` plugin manifest schema.
- The `agent.Submission`, `agent.Verdict`, `badge.Claims`,
  `federation.Signal` Go types, when consumed from `pkg/` (we do not
  guarantee internal package stability — keep your imports outside
  `internal/`).
- The CLI flags on every shipped binary (`joblantern`, `joblantern-bot`,
  every `mcp-*`).
- Migration numbering and naming (no renaming an existing migration).

## Not API-stable

- Internal Go packages (`internal/...`). Refactors land freely.
- The `risk.Bands` defaults and the exact set of `pattern.rules.yaml`
  rules — these are operator-tunable knobs; we change the defaults
  whenever the learning pipeline (Phase 34) finds a better value.
- Generated SQL bindings (`internal/db/`).
- Container image labels and the precise contents of the distroless
  base layer.

## Deprecation policy

When we plan to remove or change an API-stable surface:

1. Mark it deprecated in the relevant doc and source comment.
2. Keep the old behaviour functional for **at least one minor release**.
3. Emit a slog warning at use-site when the deprecated path is taken.
4. Remove only at the next major version.

## Browser / extension compatibility

The WebExtension targets Chrome MV3 (current Stable −2 versions) and
Firefox ESR. Older browsers may still load it but are not part of the
test matrix.

## Postgres compatibility

We test against Postgres 16. Postgres 15 is best-effort; 17 will be
added to the test matrix when it stabilises in Debian.

## Go version

Joblantern compiles with Go 1.25+. CI pins the major version; we will
not adopt a Go feature that forces operators to upgrade their build
environment without a release-note flag.
