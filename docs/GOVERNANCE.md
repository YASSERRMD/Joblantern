# Governance

Joblantern is a small project today and a public utility tomorrow.
This document is the social contract.

## Roles

| Role | Who | Powers |
|---|---|---|
| **Maintainer** | Mohamed Yasser (`@YASSERRMD`) | Merges PRs, cuts releases, decides ADRs. |
| **Reviewer** | Anyone who has merged ≥ 3 substantive PRs | May `Approve` PRs; merge requires a Maintainer. |
| **Contributor** | Anyone with a merged PR | Listed in `CONTRIBUTORS.md`. |
| **User / Operator** | NGOs and individuals running Joblantern | No formal status; their issues drive the roadmap. |

A second Maintainer is added when (a) sustained contributor presence
exists and (b) the existing Maintainer publicly nominates them via a
PR amending this document.

## Decision-making

- **Code change** — open a PR. Maintainer reviews; for substantive
  changes a second Reviewer's Approve is required.
- **Architectural change** — open an ADR PR under `docs/ADR/` first.
  Merge only after the ADR is accepted.
- **Roadmap change** — edit `docs/ROADMAP.md` in the same PR that
  introduces the new direction.
- **Disagreement** — discuss in the issue / PR. Final say rests with
  the Maintainer; that decision is reversible by a future PR.

## Communication

- Bug reports, feature requests, roadmap discussion → GitHub Issues.
- Security vulnerabilities → see `SECURITY.md` (private channel).
- General chatter / questions → GitHub Discussions if enabled,
  otherwise issues.

## Conflict of interest

Maintainers, Reviewers, and Contributors disclose any commercial
relationship with a recruiter or job board whose listing is being
verified or whose data is being integrated. Disclosure goes in the PR
description. The Maintainer recuses from merging PRs in which they
hold an interest.

## License + DCO

All contributions are licensed under Apache 2.0 and require a
`Signed-off-by` trailer (DCO 1.1). The CI lint enforces both.

## Inactive-Maintainer policy

If the Maintainer is unreachable for 90 consecutive days, the most
recent active Reviewer may open a PR amending this document to elect a
new Maintainer. The PR remains open for a 14-day public-comment window
before merge.

## This document is amendable

Open a PR. Two Approves (Maintainer + Reviewer) required.
