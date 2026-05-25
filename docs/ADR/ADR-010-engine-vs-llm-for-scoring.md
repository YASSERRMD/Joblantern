# ADR-010 — Risk score is computed by code, not by an LLM

- **Status:** Accepted
- **Date:** 2025-11-09

## Decision

The numeric verdict (`overall_risk`, `confidence`) is computed by
`internal/risk.Score`. The LLM is free to author natural-language
narrative around the result, but it never decides whether a listing
is green / yellow / red.

## Rationale

- Reproducibility: same input → same output, every time.
- Auditability: every weight applied is exposed in `WeightApplied`.
- Cost & latency: a regex + arithmetic is free; an LLM call per
  verdict is not.
- Defence in depth: a prompt-injection attempt that tries to
  influence the verdict has no path to it.
