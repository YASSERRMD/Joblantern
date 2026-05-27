# Knowledge Graph

Beyond verdicts, Joblantern publishes a queryable knowledge graph of
recruitment-fraud entities so investigators, journalists, and
researchers can ask structured questions.

## Surfaces

- `/sparql` — SPARQL 1.1 endpoint ([endpoint.go](../internal/kg/sparql/endpoint.go))
- `/kg/builder` — journalist-friendly query builder ([builder](../internal/kg/builder/journalist.go))
- `/kg/<entity>` — Linked Open Data resource ([LOD](../internal/kg/lod/contentneg.go))
- `/.well-known/void.ttl` — VoID dataset description

## Schema

Predicates and IRIs are stable. See [schema](../internal/kg/schema/mapping.go).
Vocabulary aligns with schema.org and Wikidata via [vocab](../internal/kg/vocab/align.go).

## Provenance

Every triple is shipped as a quad — the named graph names the source.
A leaked dataset cannot be repurposed as "Joblantern asserted" without
the original source attribution.

## Federation

The endpoint can answer queries that span Wikidata and DBpedia. See
[federation](../internal/kg/federation/peers.go).

## Cookbook

See [KG-COOKBOOK](KG-COOKBOOK.md).
