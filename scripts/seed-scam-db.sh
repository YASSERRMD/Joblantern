#!/usr/bin/env bash
#
# seed-scam-db.sh
# ---------------
# Loads the curated CSV in data/seed-scams.csv into the scam_reports
# table. Idempotent on (phone, email, domain) — re-runs do not duplicate.
#
# Usage:
#   DATABASE_URL=postgres://... bash scripts/seed-scam-db.sh
#
set -euo pipefail

DATABASE_URL="${DATABASE_URL:-postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable}"
CSV="${CSV:-data/seed-scams.csv}"

if ! command -v psql >/dev/null 2>&1; then
    echo "psql is required (postgres client tools)"; exit 1
fi

if [[ ! -f "$CSV" ]]; then
    echo "CSV not found: $CSV"; exit 1
fi

psql "$DATABASE_URL" <<SQL
CREATE TEMP TABLE _seed (
    company_name text, address text, lat float8, lon float8,
    phone text, email text, domain text,
    report_source text, report_url text, summary text, reported_at timestamptz
);

\\copy _seed FROM '$CSV' WITH (FORMAT csv, HEADER true);

INSERT INTO scam_reports
    (company_name, address, address_geom, phone, email, domain, report_source, report_url, summary, reported_at)
SELECT s.company_name, s.address,
       ST_SetSRID(ST_MakePoint(s.lon, s.lat), 4326)::geography,
       s.phone, s.email, s.domain,
       s.report_source, s.report_url, s.summary, s.reported_at
  FROM _seed s
 WHERE NOT EXISTS (
        SELECT 1 FROM scam_reports r
         WHERE COALESCE(r.phone,'')  = COALESCE(s.phone,'')
           AND COALESCE(r.email,'')  = COALESCE(s.email,'')
           AND COALESCE(r.domain,'') = COALESCE(s.domain,'')
       );
SQL

echo "OK: seed loaded from $CSV"
