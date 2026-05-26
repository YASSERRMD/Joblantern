-- +goose Up
-- recruiter_orgs + recruiter_badges
-- ---------------------------------
-- A recruiter or job board can register an org, get an API key, submit
-- a listing for pre-verification, and receive a signed badge they can
-- embed publicly. The badge ID is a public verifier handle.
--
-- v1 keeps this deliberately small: no pricing tiers, no member roles,
-- no manual-review KYB workflow. Those are operator concerns.

CREATE TABLE recruiter_orgs (
    id            uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          text        NOT NULL,
    contact_email citext      NOT NULL,
    api_key_hash  text        NOT NULL UNIQUE,
    -- trust_level mirrors the federation trust idiom:
    --   vouched   = manually KYB-checked by an operator
    --   observed  = self-registered, badges valid but ranked lower
    --   suspended = badges no longer verify
    trust_level   text        NOT NULL DEFAULT 'observed'
                              CHECK (trust_level IN ('vouched','observed','suspended')),
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX recruiter_orgs_trust_idx ON recruiter_orgs (trust_level);

CREATE TABLE recruiter_badges (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id          uuid        NOT NULL REFERENCES recruiter_orgs(id) ON DELETE CASCADE,
    verification_id uuid        NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,

    -- The signed badge token. Independent of the verification id so
    -- the public verifier endpoint /badge/<token> never leaks an
    -- internal verification id.
    token           text        NOT NULL UNIQUE,
    risk            text        NOT NULL CHECK (risk IN ('green','yellow','red')),
    issued_at       timestamptz NOT NULL DEFAULT now(),
    expires_at      timestamptz NOT NULL,
    revoked_at      timestamptz
);

CREATE INDEX recruiter_badges_org_idx          ON recruiter_badges (org_id);
CREATE INDEX recruiter_badges_verification_idx ON recruiter_badges (verification_id);
CREATE INDEX recruiter_badges_expires_idx      ON recruiter_badges (expires_at);

-- +goose Down
DROP TABLE IF EXISTS recruiter_badges;
DROP TABLE IF EXISTS recruiter_orgs;
