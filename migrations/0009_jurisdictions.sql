-- +goose Up
-- jurisdictions
-- -------------
-- Per-country rule pack for recruitment-fee legality and licensing. The
-- mcp-law server reads from this table (and from a YAML data file for
-- richer fields; the table is the canonical "is this jurisdiction known
-- to us yet?" surface).
--
--   code             ISO 3166-1 alpha-2 (e.g. AE, SA, PH)
--   name             human-readable country name
--   recruitment_fee_legal  is it legal to charge a job seeker a fee?
--   max_fee_pct      legal cap on the fee, as % of first-month salary
--                    (NULL = no specified cap, or fee illegal entirely)
--   citation_url     authoritative source for the above (regulator or
--                    statute). Every row MUST have one — Joblantern
--                    never tells a user "this is illegal" without
--                    pointing them at the source.
--   updated_at       audit; review cadence is documented in ADR-009.

CREATE TABLE jurisdictions (
    code                   text        PRIMARY KEY,
    name                   text        NOT NULL,
    recruitment_fee_legal  boolean     NOT NULL,
    max_fee_pct            numeric(5,2) CHECK (max_fee_pct IS NULL OR (max_fee_pct >= 0 AND max_fee_pct <= 100)),
    citation_url           text        NOT NULL,
    updated_at             timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS jurisdictions;
