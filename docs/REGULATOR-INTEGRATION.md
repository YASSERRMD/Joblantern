# Regulator Integration

Joblantern accepts feeds from verified regulators and publishes
aggregated intelligence back. The design treats regulator data as
**weighted input**, never as absolute truth — false-positive risk and
appeals are real concerns even with official designations.

## Account verification

A regulator account requires:

1. Registration from an **official domain** with a `TXT joblantern-regulator=<id>` marker.
2. A **signed letter** from a published agency contact, verified offline by the [Trust & Safety Council](GOVERNANCE.md).
3. Optional **mutual-TLS** client certificate for high-volume feeds.

## Inbound feeds

| Feed       | Weight                                 | Source            |
|------------|----------------------------------------|-------------------|
| Blacklist  | risk score floor of 95/100             | `blacklist.Entry` |
| Whitelist  | risk score reduction by 20             | `whitelist.Entry` |

Blacklist entries are subject to appeal via the Trust & Safety Council.
A successful appeal removes the floor but the regulator is notified.

## Outbound feeds

| Feed                          | Audience       | Aggregation level |
|-------------------------------|----------------|-------------------|
| National dashboard            | Regulator UI   | Country-only, k≥5 |
| Consented complaint forwarding| Regulator API  | Per-verdict, with consent |
| Signed bulletins              | Public         | As-published |

## Audit

Every regulator action is appended to a hash-chained audit log.
The tip hash is published in the Phase 9 transparency log so any
tampering is publicly detectable.

## See also

- [MOU-TEMPLATE](MOU-TEMPLATE.md) — bilateral agreement scaffold.
- [GOVERNANCE](GOVERNANCE.md) — appeals and council oversight.
- [TRANSPARENCY](TRANSPARENCY.md) — annual public reports.
