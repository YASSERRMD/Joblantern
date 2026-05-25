-- Spatial and semantic scam_reports lookups.
--
-- Used by mcp-scam-db to expose:
--   * search_reports_near_point — "how many known-bad reports are
--     within R metres of this address?"
--   * semantic_search           — "which existing reports most
--     resemble this freeform recruiter message?"

-- name: SearchScamReportsNearPoint :many
-- @lon, @lat are WGS84; @radius_m is metres. Geography type makes
-- ST_DWithin's third argument metres natively.
SELECT *,
       ST_Distance(
           address_geom,
           ST_SetSRID(ST_MakePoint(@lon::float8, @lat::float8), 4326)::geography
       ) AS distance_m
  FROM scam_reports
 WHERE address_geom IS NOT NULL
   AND ST_DWithin(
           address_geom,
           ST_SetSRID(ST_MakePoint(@lon::float8, @lat::float8), 4326)::geography,
           @radius_m::float8
       )
 ORDER BY distance_m ASC
 LIMIT @lim::int;

-- name: CountScamReportsNearPoint :one
-- Faster path when the caller only needs density, not the rows.
SELECT COUNT(*)::bigint AS n
  FROM scam_reports
 WHERE address_geom IS NOT NULL
   AND ST_DWithin(
           address_geom,
           ST_SetSRID(ST_MakePoint(@lon::float8, @lat::float8), 4326)::geography,
           @radius_m::float8
       );

-- name: SearchScamReportsByEmbedding :many
-- Cosine distance over the HNSW index. Lower distance = more similar.
-- The embedding is passed as a pgvector string literal ("[0.1,0.2,...]")
-- — sqlc treats it as text and pgvector parses it on the server side.
SELECT *,
       (embedding <=> @needle::vector) AS cosine_distance
  FROM scam_reports
 WHERE embedding IS NOT NULL
 ORDER BY embedding <=> @needle::vector
 LIMIT @lim::int;
