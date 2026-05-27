# Active Defense

Once an instance and its NGO partners reach confidence that a domain
is a scam front, they can take active steps to make life harder for
the operator: registrar abuse reports, safe-browsing submissions,
DNS sinkholes, monitor for re-registration. Phase 44 is the
human-in-the-loop scaffolding for those measures.

## Hard rule

Joblantern never auto-submits an active-defense action. Every
packet is queued via [throttle](../internal/defense/throttle/control.go)
and requires a named human approver in the case docket.

## Surfaces

- [takedown](../internal/defense/takedown/registrar.go) — registrar abuse packet.
- [safebrowsing](../internal/defense/safebrowsing/submit.go) — APWG / Google / Microsoft.
- [patterns](../internal/defense/patterns/match.go) — auto-flag listings on confirmed domains and typo-squat variants.
- [monitor](../internal/defense/monitor/respawn.go) — variant-TLD re-registration watch.
- [blocklist](../internal/defense/blocklist/publish.go) — hosts.txt, pi-hole, uBlock.
- [coordinated](../internal/defense/coordinated/group.go) — participation in cooperative groups.
- [appeals](../internal/defense/appeals/appeal.go) — false-positive remedy.
- [effectiveness](../internal/defense/effectiveness/report.go) — monthly report.

## Performance target

A confirmed scam domain has a ready-for-review takedown packet
generated in under 60 seconds. Reviewer turnaround is tracked in the
monthly effectiveness report.

## See also

- [LEGAL-REVIEW](LEGAL-REVIEW.md)
- [GOVERNANCE](GOVERNANCE.md)
