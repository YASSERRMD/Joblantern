# Recruitment Network Graph Analysis

Individual scam reports are leaves on a tree. The interesting
structure is the network — same phones across "different" companies,
shared director names, shared payment accounts. Graph analysis
reveals the operators behind dozens of fake fronts.

## Pipeline

1. [Entity resolution](../internal/graphanalysis/entity/resolve.go) — phones, emails, addresses, directors, payment instruments become canonical nodes.
2. [Edges](../internal/graphanalysis/edges/cooccurrence.go) — co-occurrence in reports drives weights.
3. [Components](../internal/graphanalysis/components/cc.go) — connected components are candidate clusters.
4. [Centrality](../internal/graphanalysis/centrality/pagerank.go) — PageRank highlights ringleaders.
5. [Drift](../internal/graphanalysis/drift/detect.go) — fast-growing components are flagged for review.

## Surfaces

- **API**: d3-compatible JSON via `/api/v1/graph/<component>`.
- **Workbench**: investigator UI at `/workbench/graph`.
- **Exports**: GraphML and Gephi CSV for offline analysis.

## Privacy

Submitter identities never enter the graph. Only entity-side nodes
(phone, email, address, director, payment) are exposed.

## Worked example

A confirmed scam company `c1` is selected in the workbench. The
component view shows 7 sibling companies — `c1..c7` — connected
through 2 phones and 1 address. Centrality identifies the address
node as the cluster's bridge.
