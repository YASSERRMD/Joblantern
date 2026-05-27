-- +goose Up
CREATE TABLE IF NOT EXISTS marketplace_verifications (
  id              uuid PRIMARY KEY,
  submitted_at    timestamptz NOT NULL DEFAULT now(),
  country         text NOT NULL,
  platform        text,
  category        text,
  listing_url     text,
  listing_hash    text NOT NULL,
  price           numeric(12,2),
  currency        text,
  contact_phone   text,
  contact_email   text,
  payment_method  text,
  shipping_method text,
  risk_score      smallint NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
  risk_band       text NOT NULL CHECK (risk_band IN ('green','yellow','red')),
  red_flags       jsonb NOT NULL DEFAULT '[]'::jsonb
);

CREATE INDEX IF NOT EXISTS marketplace_country_idx ON marketplace_verifications(country);
CREATE INDEX IF NOT EXISTS marketplace_phone_idx   ON marketplace_verifications(contact_phone);

-- +goose Down
DROP TABLE IF EXISTS marketplace_verifications;
