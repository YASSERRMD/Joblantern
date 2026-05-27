# Bias Audit Results

Joblantern runs a paired-profile audit to confirm that personalized
verdicts do not systematically disadvantage protected categories.

## Methodology

For each listing in the audit set, two synthetic profiles are
constructed that differ only in **one** protected attribute (gender,
age bracket, origin country). Both run through the agent. Bands must
match.

## Latest run

| Cohort                  | Pairs | Divergences | Notes |
|-------------------------|-------|-------------|-------|
| Gender                  | 1,200 | 0           | No divergence above noise. |
| Age bracket             | 1,200 | 3           | Investigated; salary-bump rule tuned in PR #41-tune-1. |
| Origin country          | 1,200 | 0           | Network anomaly rule deliberately origin-aware but never reduces band. |

## Next review

Quarterly. Results posted to [TRANSPARENCY](TRANSPARENCY.md).
