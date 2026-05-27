# Embassy & Consulate Pre-Departure Kiosk

Many origin-country embassies host walk-in visits before a migrant
worker travels. A simple kiosk at that visit can flag scam offers
before the worker boards a plane, preventing some of the worst
outcomes Joblantern was built to fight.

## Flow at the kiosk

1. **Scan** the job-offer QR or a barcode on the contract.
2. **Verdict** appears in 5–30 seconds.
3. **Red**: triggers the [counseling script](../internal/embassy/counseling/red.go) and prints a counselling summary.
4. **Green**: prints the [travel-readiness checklist](../internal/embassy/checklist/green.go).
5. **Yellow**: prints both, plus the closest embassy hotline.

## Components

- [UI](../internal/embassy/ui/kiosk.go) — large-text, no-typing.
- [Officer mode](../internal/embassy/officer/mode.go) — override + counseling notes.
- [ESC/POS printer](../internal/embassy/printer/escpos.go).
- [Embassy directory](../internal/embassy/directory/directory.go).
- [Offline cache](../internal/embassy/offline/cache.go).
- [Fleet management](../internal/embassy/fleet/fleet.go).
- [Anonymized analytics](../internal/embassy/analytics/aggregate.go).

## Privacy

- The kiosk **never** stores applicant names beyond the active session.
- Officer counselling notes are sealed and only accessible by the
  same embassy.
- Fleet-level analytics are aggregated and k-anonymized before they
  leave the embassy.

## See also

- [EMBASSY-KIOSK-HW](EMBASSY-KIOSK-HW.md)
- [EMBASSY-DEPLOY](EMBASSY-DEPLOY.md)
- [REGULATOR-INTEGRATION](REGULATOR-INTEGRATION.md)
