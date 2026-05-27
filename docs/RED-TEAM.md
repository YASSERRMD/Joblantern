# Red-Team Mode

Scammers will probe Joblantern's defenses. We probe them first. The
nightly red-team run is the single most important signal that the
agent is still doing its job.

## Pipeline

1. [generator](../internal/redteam/generator/synthetic.go) — average-case scams.
2. [adversarial](../internal/redteam/adversarial/evade.go) — evasive mutations.
3. [jailbreak](../internal/redteam/jailbreak/corpus.go) — prompt-injection corpus.
4. [runner](../internal/redteam/runner/nightly.go) — drives the agent stack.
5. [detection](../internal/redteam/detection/dashboard.go) — rolling rate.
6. [regression](../internal/redteam/regression/promote.go) — missed → permanent test.
7. [defense](../internal/redteam/defense/depth.go) — defense-in-depth scoring.

## Cadence

- Nightly run, results posted by 06:00 UTC.
- Monthly report from [the template](RED-TEAM-REPORT-TEMPLATE.md).
- Quarterly review by the Trust & Safety Council.

## Bounty & disclosure

- [bounty.Triage](../internal/redteam/bounty/intake.go) sets the initial response window.
- [disclosure.Default](../internal/redteam/disclosure/policy.go) sets the public policy: 48 hour ack, 7 days for critical fixes, 90 days for disclosure.

## See also

- [GOVERNANCE](GOVERNANCE.md)
- [SECURITY](../SECURITY.md)
