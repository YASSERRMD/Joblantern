# Legitimate Interest Assessment — MCP Sources

For each MCP source we run an LIA (Legitimate Interest Assessment)
that documents why processing is necessary and proportionate.

## Source rubric (per source)

1. **Identify the legitimate interest** (preventing recruitment fraud).
2. **Necessity test** — is there a less-intrusive way? If yes, use it.
3. **Balancing test** — what reasonable expectations do data subjects have?
4. **Safeguards** — what do we do to minimise the risk to subjects?

## Sources & LIA summary

| Source            | Necessity | Balancing | Safeguards |
|-------------------|-----------|-----------|------------|
| `mcp-domain`      | High      | Neutral   | Public WHOIS only; never reverse-lookups subjects. |
| `mcp-pattern`     | Medium    | Low risk  | Aggregate patterns, no per-subject extraction. |
| `mcp-registry`    | High      | Neutral   | Public company registries, official endpoints. |
| `mcp-streetview`  | Medium    | Higher    | Off by default; consent required for storage. |
| `mcp-scam-db`     | High      | Low risk  | Hashed phone numbers only, region-scoped. |
| `mcp-salary`      | Medium    | Low risk  | Open public statistics. |
| `mcp-rental-listings` | Medium | Higher  | ToS-checked aggregators only; image hashes only. |
| `mcp-accreditation` | High    | Low risk  | Government registries. |
| `mcp-cheki`       | Medium    | Low risk  | Public catalogues. |

## Review cadence

Annual, or whenever a new MCP source is added.
