# Data retention

| Table | TTL | Notes |
|---|---|---|
| verifications | 90 days | User can wipe earlier; cron deletes older rows. |
| evidence_facts | 90 days | Cascades from verifications. |
| verification_feedback | until moderated | Then archived or deleted. |
| mcp_audit_log | 30 days | Rolled to cold storage. |
| scam_reports | indefinite | Curated dataset. |
| sessions | hard expiry (`expires_at`) | Refreshed on use. |

Operators MAY override these via SQL or via a future configuration knob;
the table is the contract.
