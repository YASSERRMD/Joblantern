-- +goose Up
-- evidence_facts
-- --------------
-- Every individual piece of evidence the agent gathered while processing
-- a verification. This is the "show your work" layer that lets us audit
-- a verdict end-to-end.
--
-- A fact carries:
--   * source / tool_name — which MCP server (e.g. 'joblantern.address')
--     and which tool produced it (e.g. 'classify_building_type').
--   * fact_type — controlled vocabulary defined in migration 0007.
--   * value — the raw structured payload (jsonb so different fact
--     types can have different shapes without a schema migration).
--   * supports_risk — does this fact pull the verdict towards green,
--     yellow, or red? (Constraint enforced in migration 0007.)
--   * weight — magnitude of the signal in [0,1]. The risk engine
--     combines weight + supports_risk into the overall score.
--   * citation — public URL the user can click to verify the source.
--   * fetched_at — when we observed it.
--
-- Indexes:
--   * (verification_id, source, tool_name) for grouping by sub-agent
--     when rendering the results page.
--   * (verification_id, fact_type) for queries that pull all facts of
--     a given type (e.g. "all address-related facts").

CREATE TABLE evidence_facts (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id uuid        NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,

    source          text        NOT NULL,
    tool_name       text        NOT NULL,
    fact_type       text        NOT NULL,
    value           jsonb       NOT NULL DEFAULT '{}'::jsonb,
    supports_risk   text        NOT NULL,
    weight          numeric(4,3) NOT NULL DEFAULT 0
                                CHECK (weight >= 0 AND weight <= 1),
    citation        text,
    fetched_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX evidence_facts_verification_idx
    ON evidence_facts (verification_id);

CREATE INDEX evidence_facts_verification_source_tool_idx
    ON evidence_facts (verification_id, source, tool_name);

CREATE INDEX evidence_facts_verification_type_idx
    ON evidence_facts (verification_id, fact_type);

-- +goose Down
DROP TABLE IF EXISTS evidence_facts;
