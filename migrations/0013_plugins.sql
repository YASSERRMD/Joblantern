-- +goose Up
-- plugins
-- -------
-- Per-deployment registry of installed third-party MCP plugins.
-- The manifest itself is stored verbatim so operators can audit it
-- without re-fetching from the source URL.

CREATE TABLE plugins (
    id           uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         text        NOT NULL UNIQUE,
    version      text        NOT NULL,
    license      text        NOT NULL,
    trust        text        NOT NULL DEFAULT 'community'
                             CHECK (trust IN ('official','community','external')),
    manifest     text        NOT NULL,
    source_url   text,
    enabled      boolean     NOT NULL DEFAULT true,
    installed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX plugins_trust_idx   ON plugins (trust);
CREATE INDEX plugins_enabled_idx ON plugins (enabled);

-- +goose Down
DROP TABLE IF EXISTS plugins;
