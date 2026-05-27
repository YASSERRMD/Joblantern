-- +goose Up
-- Phase 40: education & visa-mill verdicts.
CREATE TABLE IF NOT EXISTS edu_verifications (
  id              uuid PRIMARY KEY,
  submitted_at    timestamptz NOT NULL DEFAULT now(),
  country         text NOT NULL,
  institution     text NOT NULL,
  program         text,
  agent_email     text,
  agent_phone     text,
  visa_pathway    text,
  tuition_amount  numeric(12,2),
  tuition_currency text,
  risk_score      smallint NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
  risk_band       text NOT NULL CHECK (risk_band IN ('green','yellow','red')),
  red_flags       jsonb NOT NULL DEFAULT '[]'::jsonb,
  joined_job_id   uuid REFERENCES verifications(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS edu_verifications_institution_idx ON edu_verifications(lower(institution));
CREATE INDEX IF NOT EXISTS edu_verifications_country_idx     ON edu_verifications(country);

-- +goose Down
DROP TABLE IF EXISTS edu_verifications;
