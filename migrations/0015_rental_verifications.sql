-- +goose Up
-- Phase 39: rental_verifications stores housing scam verdicts. Schema
-- mirrors verifications but keeps a separate table so retention and
-- access policies can diverge.
CREATE TABLE IF NOT EXISTS rental_verifications (
  id              uuid PRIMARY KEY,
  submitted_at    timestamptz NOT NULL DEFAULT now(),
  country         text NOT NULL,
  city            text,
  listing_url     text,
  listing_hash    text NOT NULL,
  monthly_rent    numeric(12,2),
  currency        text,
  deposit_method  text,
  contact_phone   text,
  contact_email   text,
  landlord_name   text,
  risk_score      smallint NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
  risk_band       text NOT NULL CHECK (risk_band IN ('green','yellow','red')),
  red_flags       jsonb NOT NULL DEFAULT '[]'::jsonb,
  joined_job_id   uuid REFERENCES verifications(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS rental_verifications_country_idx ON rental_verifications(country);
CREATE INDEX IF NOT EXISTS rental_verifications_phone_idx   ON rental_verifications(contact_phone);
CREATE INDEX IF NOT EXISTS rental_verifications_hash_idx    ON rental_verifications(listing_hash);

-- +goose Down
DROP TABLE IF EXISTS rental_verifications;
