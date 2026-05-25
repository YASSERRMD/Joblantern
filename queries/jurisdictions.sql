-- jurisdictions lookups for the mcp-law server.

-- name: GetJurisdiction :one
SELECT *
  FROM jurisdictions
 WHERE code = $1;

-- name: ListJurisdictions :many
SELECT *
  FROM jurisdictions
 ORDER BY code ASC;

-- name: UpsertJurisdiction :exec
-- Used by the periodic loader that refreshes the table from
-- data/recruitment-law.json (see Phase 11).
INSERT INTO jurisdictions (code, name, recruitment_fee_legal, max_fee_pct, citation_url, updated_at)
VALUES ($1, $2, $3, $4, $5, now())
ON CONFLICT (code) DO UPDATE
    SET name                  = EXCLUDED.name,
        recruitment_fee_legal = EXCLUDED.recruitment_fee_legal,
        max_fee_pct           = EXCLUDED.max_fee_pct,
        citation_url          = EXCLUDED.citation_url,
        updated_at            = now();
