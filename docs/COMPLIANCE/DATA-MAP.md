# Data Map

A line-by-line view of where personal data lives and flows.

## Inputs

| Surface             | Data class                                    | Trust boundary |
|---------------------|-----------------------------------------------|----------------|
| `/verifications`    | Listing text, recruiter contacts, applicant CV (optional) | Joblantern core |
| `/rental/submit`    | Listing text, landlord contacts, photos       | Joblantern core |
| `/edu/submit`       | Institution claim, program, agent contacts    | Joblantern core |
| `/api/v1` (research)| Anonymised verdict snapshots only             | Research surface |
| Federation peers    | Verdict hashes, no raw inputs                  | Federation surface |
| Regulator feeds     | Blacklist/whitelist rows                        | Regulator surface |

## Internal stores

| Store               | Contents                                       | Retention |
|---------------------|------------------------------------------------|-----------|
| `verifications`     | Per-verdict facts                              | 24 months default |
| `evidence`          | MCP-fetched artifacts                          | 6 months |
| `audit_log`         | Hash-chained ops events                        | 7 years |
| `cv_in_memory_only` | CV bytes (default)                             | request-scoped |

## Outputs

- Annual archive (DOI-citable). See [ARCHIVAL](../ARCHIVAL.md).
- Research API (anonymised). See [RESEARCH-API](../RESEARCH-API.md).
- Regulator dashboards (national aggregates).
- Transparency report (public).
