-- +goose Up
-- sessions
-- --------
-- Signed-cookie session records for the Joblantern web UI.
--
--   id          UUID — opaque session identifier (matches the cookie value)
--   user_id     FK → users.id (cascade on delete so logout-everywhere works)
--   token_hash  hash of the secret half of the session token, never plain
--   expires_at  hard expiry; sessions older than this are rejected
--   created_at  audit
--
-- The cookie sent to the browser is split: an id half (looked up in this
-- table) and a secret half (hashed and compared on each request). This
-- prevents a leaked DB row from being directly usable as a session.

CREATE TABLE sessions (
    id          uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  text        NOT NULL,
    expires_at  timestamptz NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx       ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx    ON sessions (expires_at);

-- +goose Down
DROP TABLE IF EXISTS sessions;
