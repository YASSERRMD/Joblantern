-- CRUD for the verifications table. Consumed by:
--   * cmd/joblantern's web layer (submit, fetch, list)
--   * internal/agent (read while running, update on completion)
--
-- claimed_geom is written via ST_SetSRID(ST_MakePoint(lon,lat),4326)::geography
-- in the query rather than asking sqlc to model PostGIS types. NULL is
-- written when the agent has not yet geocoded the address.

-- name: CreateVerification :one
INSERT INTO verifications (
    user_id,
    raw_input,
    listing_url,
    recruiter_email,
    recruiter_phone,
    company_name,
    claimed_address,
    claimed_salary,
    role,
    jurisdiction
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetVerification :one
SELECT *
  FROM verifications
 WHERE id = $1;

-- name: ListVerificationsByUser :many
SELECT *
  FROM verifications
 WHERE user_id = $1
 ORDER BY created_at DESC
 LIMIT $2
 OFFSET $3;

-- name: ListRecentVerifications :many
SELECT *
  FROM verifications
 ORDER BY created_at DESC
 LIMIT $1;

-- name: SetVerificationGeom :exec
UPDATE verifications
   SET claimed_geom = ST_SetSRID(ST_MakePoint(@lon::float8, @lat::float8), 4326)::geography
 WHERE id = @id;

-- name: SetVerificationRunning :exec
UPDATE verifications
   SET status = 'running'
 WHERE id = $1
   AND status = 'pending';

-- name: CompleteVerification :exec
UPDATE verifications
   SET status        = 'completed',
       overall_risk  = $2,
       confidence    = $3,
       completed_at  = now()
 WHERE id = $1;

-- name: FailVerification :exec
UPDATE verifications
   SET status       = 'failed',
       completed_at = now()
 WHERE id = $1;

-- name: DeleteVerification :exec
DELETE FROM verifications
 WHERE id = $1;
