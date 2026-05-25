-- Non-spatial scam_reports lookups (phone, email, domain, name).
-- Spatial and semantic searches live in queries/scam_reports_spatial.sql.

-- name: SearchScamReportsByPhone :many
-- Phones are normalised to E.164 by the caller before insert and before
-- search; the equality check here is exact.
SELECT *
  FROM scam_reports
 WHERE phone = $1
 ORDER BY reported_at DESC NULLS LAST
 LIMIT $2;

-- name: SearchScamReportsByEmail :many
SELECT *
  FROM scam_reports
 WHERE email = $1
 ORDER BY reported_at DESC NULLS LAST
 LIMIT $2;

-- name: SearchScamReportsByEmailDomain :many
-- Match all reports whose email belongs to a given domain (e.g. catch
-- "joe@acme.test" and "abuse@acme.test" with input "acme.test").
SELECT *
  FROM scam_reports
 WHERE email ILIKE '%@' || $1
 ORDER BY reported_at DESC NULLS LAST
 LIMIT $2;

-- name: SearchScamReportsByDomain :many
-- Match the exact domain OR any subdomain of it.
SELECT *
  FROM scam_reports
 WHERE domain = $1
    OR domain ILIKE '%.' || $1
 ORDER BY reported_at DESC NULLS LAST
 LIMIT $2;

-- name: SearchScamReportsByCompanyName :many
-- pg_trgm similarity. Threshold is provided by the caller so the agent
-- can sweep from strict (0.6) to loose (0.3) progressively.
SELECT *,
       similarity(company_name, @needle::text) AS sim
  FROM scam_reports
 WHERE company_name % @needle::text
   AND similarity(company_name, @needle::text) >= @min_sim::float4
 ORDER BY sim DESC
 LIMIT @lim::int;
