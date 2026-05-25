# MCP Server — `joblantern.pattern`

## Purpose

Deterministic, rule-based red-flag classifier for recruitment listings.
Rules live in `internal/pattern/rules.yaml` and are hot-reloadable in
future versions.

## Tools

| Name | Purpose |
|---|---|
| `analyze_listing_text` | Returns matched red flags + composite score [0,1]. |
| `detect_red_flag_phrases` | Same matches without scoring. |
| `language_mismatch_check` | Flags Cyrillic/CJK in Latin-script jurisdictions. |

### Rule codes (v1)

`upfront_fee`, `extraordinary_pay`, `no_experience_needed`,
`urgency`, `untraceable_contact`, `identity_documents_upfront`,
`fake_government_endorsement`.

### Error codes

`INVALID_ARGS`, `UNSUPPORTED_LANGUAGE`, `EMBEDDING_UNAVAILABLE`.
