# Long-Term Archival & Historical Record

Recruitment-fraud history is itself an important historical record.
Joblantern publishes annual archives so a researcher in 2045 can
cite a 2026 verdict with confidence.

## Annual archives

- Format: Parquet + JSON-LD ([format](../internal/archival/format/parquet.go)).
- Integrity: SHA-256 per file plus a manifest hash ([integrity](../internal/archival/integrity/hash.go)).
- DOI: assigned through DataCite ([doi](../internal/archival/doi/datacite.go)).
- Landing: BibTeX-shipping HTML page ([landing](../internal/archival/landing/page.go)).
- Mirrors: Internet Archive + Zenodo ([thirdparty](../internal/archival/thirdparty/mirror.go)).

## Anonymisation audit

Every archive runs through [anonaudit](../internal/archival/anonaudit/audit.go)
before publication. Fatal checks block release.

## Time machine

The [time machine](../internal/archival/timemachine/asof.go) lets a
user reconstruct the state of a verdict as of any past date, using
the matching annual archive plus monthly deltas.

## Survivability

- [ARCHIVAL-RETENTION](ARCHIVAL-RETENTION.md)
- [ARCHIVAL-ESCROW](ARCHIVAL-ESCROW.md)

## Citing the archive

```
Joblantern Annual Verdict Archive, 2026,
doi:10.<prefix>/joblantern.2026.
```
