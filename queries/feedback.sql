-- name: InsertVerificationFeedback :one
INSERT INTO verification_feedback (verification_id, outcome, comment)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CountFeedbackByOutcome :many
SELECT outcome, COUNT(*)::bigint AS n
  FROM verification_feedback
 GROUP BY outcome
 ORDER BY outcome;

-- name: ListPendingScamFeedback :many
-- Pending = outcome='confirmed_scam' and not yet promoted into scam_reports.
SELECT vf.id, vf.verification_id, vf.comment, vf.submitted_at,
       v.company_name, v.claimed_address, v.recruiter_email, v.recruiter_phone
  FROM verification_feedback vf
  JOIN verifications v ON v.id = vf.verification_id
 WHERE vf.outcome = 'confirmed_scam'
 ORDER BY vf.submitted_at ASC
 LIMIT $1;
