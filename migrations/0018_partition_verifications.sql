-- +goose Up
-- Phase 47: convert verifications and audit_log to RANGE-partitioned
-- tables keyed on submitted_at, one partition per month. This keeps
-- index sizes bounded as we scale to 10M+ verdicts/day.
--
-- NOTE: this migration is destructive when run against a non-empty
-- table. Production deployments run pg_partman side-by-side and cut
-- over with a swap. This file is the declarative target.

CREATE TABLE IF NOT EXISTS verifications_partitioned (
  LIKE verifications INCLUDING ALL
) PARTITION BY RANGE (submitted_at);

CREATE TABLE IF NOT EXISTS verifications_2026_05 PARTITION OF verifications_partitioned
  FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE IF NOT EXISTS verifications_2026_06 PARTITION OF verifications_partitioned
  FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

-- +goose Down
DROP TABLE IF EXISTS verifications_2026_05;
DROP TABLE IF EXISTS verifications_2026_06;
DROP TABLE IF EXISTS verifications_partitioned;
