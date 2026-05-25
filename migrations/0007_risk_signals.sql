-- +goose Up
-- Tighten the evidence_facts.supports_risk and .fact_type columns by
-- replacing free-form text with controlled vocabularies. We use CHECK
-- constraints rather than native ENUM types so adding a value later is
-- a simple ALTER TABLE … DROP CONSTRAINT … ADD CONSTRAINT, with no
-- ENUM-rebuild gymnastics.
--
-- supports_risk vocabulary:
--   green   the fact reduces overall risk (e.g. registered company,
--           old domain, salary in line)
--   yellow  the fact is ambiguous and warrants more scrutiny
--   red     the fact increases overall risk (e.g. address is a
--           residential building, recruiter demands upfront fee)
--   neutral the fact is informational only (e.g. reverse-geocoded
--           administrative area)
--
-- fact_type vocabulary — see docs/RISK-MODEL.md for the canonical list.
-- We seed enough types to cover Phases 04-12; new types are added by
-- future migrations alongside the MCP server that introduces them.

ALTER TABLE evidence_facts
    ADD CONSTRAINT evidence_facts_supports_risk_check
    CHECK (supports_risk IN ('green','yellow','red','neutral'));

ALTER TABLE evidence_facts
    ADD CONSTRAINT evidence_facts_fact_type_check
    CHECK (fact_type IN (
        -- address / streetview
        'address.exists',
        'address.not_found',
        'address.building_type',
        'address.residential_match',
        'address.commercial_match',
        'address.cluster_density',
        'streetview.imagery_present',
        'streetview.imagery_absent',
        'streetview.imagery_age',
        -- registry / domain
        'registry.company_found',
        'registry.company_missing',
        'registry.company_age',
        'registry.address_match',
        'registry.address_mismatch',
        'registry.director_count',
        'domain.age',
        'domain.ssl_first_seen',
        'domain.archive_first_snapshot',
        'domain.freshness_score',
        -- pattern / scam DB
        'pattern.red_flag',
        'pattern.language_mismatch',
        'pattern.similarity_to_reports',
        'scamdb.phone_match',
        'scamdb.email_match',
        'scamdb.domain_match',
        'scamdb.name_fuzzy_match',
        'scamdb.spatial_cluster',
        -- salary / law
        'salary.within_range',
        'salary.implausibly_high',
        'salary.implausibly_low',
        'law.recruitment_fee_illegal',
        'law.recruitment_fee_legal',
        'law.licensing_required',
        -- routing
        'routing.reachable',
        'routing.unreachable',
        'routing.commute_implausible'
    ));

-- +goose Down
ALTER TABLE evidence_facts DROP CONSTRAINT IF EXISTS evidence_facts_fact_type_check;
ALTER TABLE evidence_facts DROP CONSTRAINT IF EXISTS evidence_facts_supports_risk_check;
