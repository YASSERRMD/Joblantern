-- +goose Up
-- scam_reports
-- ------------
-- The local catalogue of known-bad recruiters / addresses / phones /
-- emails / domains, curated from public reports and (later, with
-- moderation) from user feedback.
--
-- Lookups happen from the mcp-scam-db server using:
--   * exact / normalised phone match
--   * exact email + domain-only match
--   * exact + subdomain domain match
--   * pg_trgm fuzzy company-name match
--   * ST_DWithin spatial proximity around address_geom
--   * HNSW cosine similarity over `embedding` (semantic / freeform text)
--
-- Embedding dimensionality is 384 (matches nomic-embed-text / many small
-- open embedding models). Migrations can introduce additional embedding
-- columns of other dimensions in future without breaking this one.

CREATE TABLE scam_reports (
    id              uuid        PRIMARY KEY DEFAULT uuid_generate_v4(),

    company_name    text,
    address         text,
    address_geom    geography(Point, 4326),
    phone           text,
    email           text,
    domain          text,

    report_source   text,
    report_url      text,
    summary         text,
    embedding       vector(384),

    reported_at     timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now()
);

-- Spatial lookup.
CREATE INDEX scam_reports_address_geom_idx
    ON scam_reports USING gist (address_geom);

-- Fuzzy company-name lookup.
CREATE INDEX scam_reports_company_name_trgm_idx
    ON scam_reports USING gin (company_name gin_trgm_ops);

-- Exact / prefix lookups.
CREATE INDEX scam_reports_phone_idx   ON scam_reports (phone);
CREATE INDEX scam_reports_email_idx   ON scam_reports (email);
CREATE INDEX scam_reports_domain_idx  ON scam_reports (domain);

-- Semantic similarity via HNSW (cosine).
-- pgvector HNSW parameters: m=16, ef_construction=64 are the upstream
-- defaults and are fine for v1 scale.
CREATE INDEX scam_reports_embedding_hnsw_idx
    ON scam_reports USING hnsw (embedding vector_cosine_ops);

-- +goose Down
DROP TABLE IF EXISTS scam_reports;
