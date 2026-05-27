# Article 30 — Records of Processing Activities (RoPA)

> Maintained per GDPR Article 30. Reviewed quarterly by the DPO.

## Controller

- Joblantern Foundation
- DPO contact: dpo@joblantern.org
- Establishment: see [DATA-RESIDENCY](../DATA-RESIDENCY.md)

## Processing activities

### 1. Verdict generation

- Purpose: detect recruitment fraud
- Categories: applicants (data subjects), recruiters (referenced entities)
- Data: listing text, contacts, optional CV
- Recipients: NGO partners (anonymised), regulators (consented), public (anonymised stats)
- Cross-border: per regional tenant; no inter-region transfer
- Retention: 24 months default, 6 months for evidence artifacts
- Security: hash-chained audit, RLS, mTLS to regulator feeds

### 2. Researcher API

- Purpose: support studies on labor trafficking
- Categories: anonymised verdicts only
- Recipients: vetted researchers, journalists, academics
- Retention: indefinite for anonymised archives

### 3. Multi-tenant operations

- Purpose: hosted-service tenant management
- Categories: tenant admin contacts only
- Retention: per tenant contract, default 24 months after offboarding
