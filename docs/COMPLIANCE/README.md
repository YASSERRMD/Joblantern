# Compliance Artifacts

This directory collects the operational artifacts a Joblantern instance
needs to demonstrate lawful operation under GDPR, India DPDPA, and
UAE/Saudi PDPL.

## Contents

- [DATA-MAP](DATA-MAP.md)
- [DPIA-TEMPLATE](DPIA-TEMPLATE.md)
- [ARTICLE-30-REGISTER](ARTICLE-30-REGISTER.md)
- [LIA-MCP-SOURCES](LIA-MCP-SOURCES.md)
- [PRIVACY-POLICY-LANGS](PRIVACY-POLICY-LANGS.md)
- [BREACH-RUNBOOK](BREACH-RUNBOOK.md)
- [DPA-TEMPLATE](DPA-TEMPLATE.md)

## Code

- [sar](../../internal/compliance/sar/request.go) — Subject Access Request workflow.
- [cookies](../../internal/compliance/cookies/banner.go) — consent banner.
- [dpdpa](../../internal/compliance/dpdpa/consent.go) — India consent flow.
- [pdpl](../../internal/compliance/pdpl/residency.go) — UAE/Saudi residency.
- [audit](../../internal/compliance/audit/iso27001.go) — ISO 27001 retention.

## Public commitments

- A deletion request is verifiably executed within 30 days.
- Cookies beyond strictly-necessary require explicit consent.
- Cross-region data transfer requires an approved mechanism.
- A material policy change is announced ≥30 days in advance.
