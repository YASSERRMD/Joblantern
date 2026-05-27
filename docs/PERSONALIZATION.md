# Skills-Based Verdict Personalization

The same listing means different things to different people. A
welder relocating to Saudi Arabia at construction-industry wages is
normal; a software engineer relocating to the same role and salary
is suspicious. Phase 41 personalizes verdicts when the user opts in
by uploading a CV.

## Pipeline

1. **Upload** — optional, max 5 MiB, multipart form. See [cv/upload.go](../internal/personalization/cv/upload.go).
2. **Parse** — local model extracts role, years of experience, skills, origin country.
3. **Compare** — [fit](../internal/personalization/fit/score.go), [salary](../internal/personalization/salary/tier.go), [location](../internal/personalization/location/anomaly.go).
4. **Explain** — [explain.Render](../internal/personalization/explain/reason.go) produces the user-facing reason.
5. **Discard** — CV bytes are zeroed unless the user opts in to retain.

## Privacy contract

- Default: **nothing persists**. CV bytes live only in memory.
- Opt-in `research-only`: aggregated, anonymised features only.
- Opt-in `improve-agent`: structured CV stored under the user's
  account; can be deleted at any time via the Phase 48 SAR flow.

## Bias audit

The audit harness in [bias](../internal/personalization/bias/audit.go)
runs the same listing across paired profiles that differ only in a
protected attribute. Divergences are surfaced to the Trust & Safety
council monthly. Results are published — see
[bias-audit](BIAS-AUDIT.md).

## Try-as-different-profile

The [simulate](../internal/personalization/simulate/profiles.go) pack
ships five demo personas so journalists and NGO trainers can show
why personalization matters without exposing real CVs.

## See also

- [BIAS-AUDIT](BIAS-AUDIT.md)
- [PRIVACY](PRIVACY.md)
