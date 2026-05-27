// Package graphql is the read-only GraphQL endpoint surface for
// vetted researchers. The schema is intentionally minimal and only
// exposes anonymized verdicts and aggregate statistics.
//
// In production this is wired through gqlgen. The schema below is the
// canonical source of truth and is shipped to the generator. Schema
// types that include personally identifying fields MUST be filtered
// before they reach the resolver layer.
package graphql

// Schema is the gqlgen schema for the public research API.
const Schema = `
"An anonymized verdict released for research."
type Verdict {
  id: ID!
  createdAt: String!
  country: String!
  industry: String!
  riskScore: Int!
  riskBand: RiskBand!
  redFlags: [String!]!
}

enum RiskBand { GREEN YELLOW RED }

type AggregateBucket {
  key: String!
  count: Int!
}

type Query {
  "Anonymized verdicts. Cursor-paginated, capped at 1000 per page."
  verdicts(country: String, since: String, after: String, first: Int = 100): [Verdict!]!
  "Counts of verdicts grouped by country."
  verdictsByCountry(since: String): [AggregateBucket!]!
  "Counts of verdicts grouped by industry."
  verdictsByIndustry(since: String): [AggregateBucket!]!
}
`
