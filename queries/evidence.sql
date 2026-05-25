-- Evidence facts insert + lookup.
--
-- The agent writes each fact as it collects it (during the parallel
-- fan-out). The web layer reads facts grouped by sub-agent (source +
-- tool) to render the results page. The risk engine reads them all
-- and computes the verdict.

-- name: InsertEvidenceFact :one
INSERT INTO evidence_facts (
    verification_id,
    source,
    tool_name,
    fact_type,
    value,
    supports_risk,
    weight,
    citation
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListEvidenceByVerification :many
SELECT *
  FROM evidence_facts
 WHERE verification_id = $1
 ORDER BY source ASC,
          tool_name ASC,
          fetched_at ASC;

-- name: ListEvidenceByVerificationAndType :many
SELECT *
  FROM evidence_facts
 WHERE verification_id = $1
   AND fact_type       = $2
 ORDER BY fetched_at ASC;

-- name: CountEvidenceByRisk :many
SELECT supports_risk, COUNT(*)::bigint AS n
  FROM evidence_facts
 WHERE verification_id = $1
 GROUP BY supports_risk
 ORDER BY supports_risk;

-- name: DeleteEvidenceForVerification :exec
DELETE FROM evidence_facts
 WHERE verification_id = $1;
