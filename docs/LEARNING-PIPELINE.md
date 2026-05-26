# Continuous learning pipeline

Joblantern does **not** retrain foundation models. The learning loop
is intentionally conservative:

1. Export anonymised `(verification, pattern_codes, outcome)` triples
   to CSV for offline analysis (operator + researcher review).
2. Compute per-rule effectiveness — precision against
   `confirmed_scam` feedback when both `confirmed_scam` and
   `confirmed_legit` are present.
3. Emit a per-rule recommendation: `keep`, `review`, or
   `consider_removal`.
4. A **human** reviews the recommendations and proposes a PR against
   `internal/pattern/rules.yaml`. The pipeline never auto-applies
   changes.

## API

```go
rows := []learning.Labelled{ ... }
learning.ExportCSV(w, rows)            // anonymised CSV
scores := learning.Effectiveness(rows) // []EffectivenessScore
fmt.Println(learning.SummaryText(scores))
```

## Thresholds

A rule needs at least `MinEvidence` (30) firings before its precision
is allowed to recommend anything other than `keep`. Below that we do
not have enough signal.

| Precision | Recommendation |
|---|---|
| ≥ 0.7 | `keep` |
| 0.4–0.7 | `review` |
| < 0.4 | `consider_removal` |

## What's deferred

- A nightly cron that runs the export + scoring and opens a PR.
- A bias audit that breaks the report down by country / language to
  catch rules that work well in one region and poorly in another.
- Drift detection alerts when the verdict distribution shifts.
- Counterfactual evaluator that replays historical verdicts against
  a proposed new rule pack before merge.
- Model-card generation for the embedding model (Phase 8 scam-DB
  semantic search).

Each of those is independently shippable on top of the primitives in
`internal/learning`.
