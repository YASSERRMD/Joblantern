# Knowledge Graph — Query Cookbook

A starter set of SPARQL queries against the Joblantern KG. Use them
as-is or as recipes to adapt.

## 1. Confirmed scam companies in a country in the last 90 days

```sparql
PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT ?company ?industry ?at WHERE {
  ?company jl:hasCountry "PH" ;
           jl:hasRiskBand "red" ;
           jl:hasIndustry ?industry ;
           jl:verifiedAt ?at .
  FILTER (?at > "2026-02-26T00:00:00"^^xsd:dateTime)
}
ORDER BY DESC(?at)
LIMIT 200
```

## 2. Companies in Manila sharing a phone with any company in Doha

```sparql
PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT ?a ?b ?phone WHERE {
  ?a jl:hasCountry "PH" ; jl:hasPhone ?phone .
  ?b jl:hasCountry "QA" ; jl:hasPhone ?phone .
  FILTER (?a != ?b)
}
```

## 3. Top 20 ringleader phones (by appearance across red verdicts)

```sparql
PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT ?phone (COUNT(?c) AS ?hits) WHERE {
  ?c jl:hasRiskBand "red" ; jl:hasPhone ?phone .
}
GROUP BY ?phone ORDER BY DESC(?hits) LIMIT 20
```

## 4. Newly-registered scam clusters this week

```sparql
PREFIX jl: <https://joblantern.org/kg/predicate/>
SELECT DISTINCT ?addr WHERE {
  ?c jl:hasAddress ?addr ; jl:hasRiskBand "red" ; jl:verifiedAt ?at .
  FILTER (?at > NOW() - "P7D"^^xsd:duration)
}
```
