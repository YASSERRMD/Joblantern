# ADR-011 — Feedback privacy model

- **Status:** Accepted
- **Date:** 2025-11-09

## Decision

User-submitted feedback rows carry only:

- The verification id.
- An outcome enum (`confirmed_scam` / `confirmed_legit` / `unsure`).
- An optional free-text comment.

Identifying information is **not** captured. Feedback only enters the
public `scam_reports` table after a human moderation step that
explicitly approves the row.

## Consequences

- Joblantern can show a useful "X% of users agreed" metric without
  ever attributing feedback to a specific person.
- Promotion to `scam_reports` is deliberate and reviewable.
