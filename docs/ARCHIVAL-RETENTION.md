# Archival Retention Policy

Joblantern annual archives are intended to survive organisational
change. The policy is binding on any successor entity that takes
over Joblantern's stewardship.

## Retention horizon

- Annual archives: indefinite.
- Monthly deltas: 7 years.
- Operational logs: per [COMPLIANCE](COMPLIANCE/README.md).

## Survivability commitments

- Mirror to at least two independent third-party archives (Internet Archive + Zenodo).
- Maintain a DataCite DOI for each annual archive.
- Hold a legal-escrow arrangement so archives outlive any single corporate form.

## Format obsolescence

Parquet and JSON-LD will be re-encoded into the then-current durable
formats every 15 years. Each re-encoding is itself archived under a
new sub-DOI; the original encoding stays available.

## What happens if Joblantern.org sunsets

The successor entity nominated in the [escrow agreement](ARCHIVAL-ESCROW.md)
takes custody of the archive copies and the DOI prefix.
