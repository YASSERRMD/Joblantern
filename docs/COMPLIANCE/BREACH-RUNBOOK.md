# Data Breach Response Runbook

## Severity classification (within 1 hour of detection)

| Sev | Definition                              | Notification |
|-----|------------------------------------------|--------------|
| 0   | Confirmed exfiltration of subject data   | 72h to authorities, immediate to council |
| 1   | Likely exfiltration                      | 72h to authorities |
| 2   | Internal exposure, no exfiltration       | Council, no external |
| 3   | Suspected (still investigating)          | Internal only |

## First 60 minutes

1. Confirm the breach is real (not a false-positive alert).
2. Stop the bleeding: revoke compromised credentials, pull the affected nodes.
3. Open the incident in the Phase 51 case docket.
4. Page the on-call: DPO, Council chair, Security lead.
5. Preserve evidence: snapshot logs, write-protect storage.

## 72 hours

- Notify the lead supervisory authority (EU: per Article 33).
- If high risk to subjects: notify data subjects (Article 34).
- Brief partner NGOs.

## Post-incident

- Public postmortem within 14 days (per [TRANSPARENCY](../TRANSPARENCY.md)).
- Add a regression test or hardening change tied to the root cause.
- Update the threat model.
