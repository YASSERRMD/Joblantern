# ADR-008 — Rules-based pattern classifier for v1

- **Status:** Accepted
- **Date:** 2025-11-09

## Decision

The pattern MCP server uses a deterministic regex rule pack
(`internal/pattern/rules.yaml`) rather than calling an LLM for
text classification.

## Rationale

- **Cost.** Listings are short; an LLM call per submission is
  unnecessary spend for a non-profit-friendly deployment.
- **Determinism.** Same input → same output. Critical for audit and
  reproducibility.
- **Latency.** Rules run in microseconds.
- **Reviewability.** Every rule is a YAML row a reviewer can read.

## Consequences

- Recall is limited to phrases we have anticipated. Mitigated by:
  - Adding rules over time.
  - Composing with `mcp-scam-db.semantic_search` (Phase 8, embedding
    similarity to existing reports) when the agent has free latency.
