-- +goose Up
-- federated_peers
-- ---------------
-- Joblantern instances that this deployment trusts as evidence peers.
-- Outgoing scam-signal exchanges flow to peers; incoming signals are
-- attributed to the peer and weighted by trust_level.
--
--   url         base URL of the peer Joblantern instance
--   name        operator-visible label
--   pubkey      ed25519 public key in raw hex; used to verify peer
--               manifests and signed scam-signal payloads
--   trust_level "vouched" | "observed" | "untrusted" — weights the
--               evidence the peer contributes; untrusted peers are
--               kept only for visibility
--   last_seen   wall-clock timestamp of the last successful handshake

CREATE TABLE federated_peers (
    id          uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    url         text        NOT NULL UNIQUE,
    name        text        NOT NULL,
    pubkey      text        NOT NULL,
    trust_level text        NOT NULL DEFAULT 'observed'
                            CHECK (trust_level IN ('vouched','observed','untrusted')),
    last_seen   timestamptz,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX federated_peers_trust_idx ON federated_peers (trust_level);

-- +goose Down
DROP TABLE IF EXISTS federated_peers;
