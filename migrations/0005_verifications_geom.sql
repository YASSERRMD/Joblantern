-- +goose Up
-- Add the resolved coordinate for the claimed address to verifications.
--
-- Why a separate migration:
--   * Splitting the geom column from the table definition keeps each
--     migration small and focused, and lets us roll back the GIS-specific
--     work without dropping the whole table.
--   * Using `geography(Point, 4326)` (WGS84 lon/lat) gives us metric
--     distances out of the box via ST_Distance / ST_DWithin without
--     extra projection logic.
--
-- GiST index lets ST_DWithin and bounding-box queries stay fast.

ALTER TABLE verifications
    ADD COLUMN claimed_geom geography(Point, 4326);

CREATE INDEX verifications_claimed_geom_idx
    ON verifications USING gist (claimed_geom);

-- +goose Down
DROP INDEX IF EXISTS verifications_claimed_geom_idx;
ALTER TABLE verifications DROP COLUMN IF EXISTS claimed_geom;
