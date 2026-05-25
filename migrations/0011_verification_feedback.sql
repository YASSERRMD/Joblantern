-- +goose Up
-- verification_feedback
-- ---------------------
-- User-confirmation feedback collected on the results page. A
-- thumbs-up/down + optional comment is stored here; a later moderation
-- step (out of v1 scope) promotes confirmed-scam feedback into
-- scam_reports rows.

CREATE TABLE verification_feedback (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id uuid        NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    outcome         text        NOT NULL CHECK (outcome IN ('confirmed_scam','confirmed_legit','unsure')),
    comment         text,
    submitted_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX verification_feedback_verification_idx
    ON verification_feedback (verification_id);
CREATE INDEX verification_feedback_outcome_idx
    ON verification_feedback (outcome);

-- +goose Down
DROP TABLE IF EXISTS verification_feedback;
