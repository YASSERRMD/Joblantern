-- +goose Up
-- verifications
-- -------------
-- One row per job-listing verification request. This is the central table
-- of Joblantern; every sub-agent run, every MCP call, every evidence
-- fact, every feedback entry hangs off a verification id.
--
-- user_id is NULLABLE because anonymous submissions are supported.
--
-- raw_input is the unprocessed payload as the user submitted it (URL,
-- recruiter message, screenshots-as-text, etc). Stored as jsonb so we
-- can evolve the submission schema without a migration.
--
-- The "claimed_*" columns capture what the listing *says* — they are
-- inputs to verification, not facts. Evidence facts live in their own
-- table (migration 0006).
--
-- status moves from 'pending' → 'running' → 'completed' (or 'failed').
-- overall_risk and confidence are populated by the risk engine on
-- completion.
--
-- jurisdiction is the destination country's ISO 3166-1 alpha-2 code
-- (e.g. AE, SA, PH). NULL means "unknown".

CREATE TABLE verifications (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         uuid        REFERENCES users(id) ON DELETE SET NULL,

    raw_input       jsonb       NOT NULL DEFAULT '{}'::jsonb,
    listing_url     text,
    recruiter_email text,
    recruiter_phone text,
    company_name    text,
    claimed_address text,
    claimed_salary  numeric(14,2),
    role            text,
    jurisdiction    text,

    status          text        NOT NULL DEFAULT 'pending'
                                CHECK (status IN ('pending','running','completed','failed')),
    overall_risk    text        CHECK (overall_risk IS NULL OR overall_risk IN ('green','yellow','red')),
    confidence      numeric(4,3) CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 1)),

    created_at      timestamptz NOT NULL DEFAULT now(),
    completed_at    timestamptz
);

CREATE INDEX verifications_user_id_idx       ON verifications (user_id);
CREATE INDEX verifications_status_idx        ON verifications (status);
CREATE INDEX verifications_created_at_idx    ON verifications (created_at DESC);
CREATE INDEX verifications_jurisdiction_idx  ON verifications (jurisdiction);

-- +goose Down
DROP TABLE IF EXISTS verifications;
