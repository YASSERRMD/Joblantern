# Joblantern Hosted Service

For NGOs that cannot run their own instance, joblantern.org offers a
hosted multi-tenant service. Tenants are isolated via row-level
security by default and can opt up to schema-per-tenant.

## Tiers

| Tier      | Audience               | Verdicts/mo | Notes |
|-----------|------------------------|-------------|-------|
| NGO Free  | Low-income NGOs        | 1,000       | Self-serve, lite KYC. |
| Pro       | Commercial recruiters  | 25,000      | Full KYC, paid. |
| Custom    | Regulators & federations | Unlimited | SLA, residency. |

## Residency

EU, US, APAC. Data never crosses a region. See [DATA-RESIDENCY](DATA-RESIDENCY.md).

## Isolation

- **rls** (default): one schema, `tenant_id` predicate enforced by Postgres RLS.
- **schema**: one schema per tenant; reserved for high-isolation customers.

## Offboarding

Tenants can pull a full archive and trigger deletion at any time.
Deletion is verifiable within 30 days under the Phase 48 SAR flow.

## Abuse

Cross-tenant enumeration is detected by passive analysis. Confirmed
abuse is escalated to the Trust & Safety Council.

## See also

- [DATA-RESIDENCY](DATA-RESIDENCY.md)
- [GOVERNANCE](GOVERNANCE.md)
- [FUNDING](FUNDING.md)
