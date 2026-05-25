-- +goose Up
-- users
-- -----
-- The authenticated user of the Joblantern web UI and (eventually) the API.
--
--   id            UUID primary key (v4)
--   email         case-insensitive unique (citext)
--   password_hash bcrypt/argon2id digest — never the plain password
--   created_at    set by default
--
-- We deliberately keep this table minimal in v1. PII collection has a
-- price; we will not store anything we do not strictly need.

CREATE TABLE users (
    id            uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         citext      NOT NULL UNIQUE,
    password_hash text        NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;
