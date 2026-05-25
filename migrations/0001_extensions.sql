-- +goose Up
-- Enable the Postgres extensions Joblantern depends on.
--
--   postgis     spatial types, distance, contains, intersects, GiST indexes
--   vector      pgvector — embedding storage and HNSW indexing
--   pg_trgm     trigram fuzzy matching (company names, addresses)
--   uuid-ossp   UUID generators (uuid_generate_v4 etc.)
--   citext      case-insensitive text columns (emails, codes)
--
-- All five ship with the joblantern/postgres image built in deploy/postgres/.

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

-- +goose Down
-- Dropping extensions while objects depend on them would fail. We only
-- drop what is safe to drop on a teardown. PostGIS is intentionally
-- left in place because it is expensive to install.

DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS vector;
-- DROP EXTENSION IF EXISTS postgis;  -- left intentionally
