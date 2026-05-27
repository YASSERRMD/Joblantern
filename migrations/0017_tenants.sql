-- +goose Up
-- Phase 45: multi-tenant hosted service. We choose row-level security
-- backed by a tenant_id column on every relevant table, plus an
-- optional schema-per-tenant for high-isolation customers.
CREATE TABLE IF NOT EXISTS tenants (
  id              uuid PRIMARY KEY,
  slug            text NOT NULL UNIQUE,
  display_name    text NOT NULL,
  tier            text NOT NULL CHECK (tier IN ('free','pro','custom')),
  region          text NOT NULL CHECK (region IN ('eu','us','apac')),
  isolation       text NOT NULL CHECK (isolation IN ('rls','schema')) DEFAULT 'rls',
  schema_name     text,
  created_at      timestamptz NOT NULL DEFAULT now(),
  suspended_at    timestamptz,
  offboarded_at   timestamptz
);

ALTER TABLE verifications ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS verifications_tenant_idx ON verifications(tenant_id);

-- +goose Down
ALTER TABLE verifications DROP COLUMN IF EXISTS tenant_id;
DROP TABLE IF EXISTS tenants;
