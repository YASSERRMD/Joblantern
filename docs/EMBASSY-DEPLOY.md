# Embassy Kiosk — Deployment Guide

> Companion to [EMBASSY-KIOSK-HW](EMBASSY-KIOSK-HW.md).

## 30-minute install

1. Flash the supplied OS image to the SBC SD card (`embassy-kiosk-${version}.img`).
2. Boot. The kiosk auto-launches into the Phase 38 UI.
3. Use a staff badge to enter Officer mode and run **Setup → Embassy registration**:
   - Country and embassy name (free-text).
   - Joblantern instance URL (defaults to your federation's hub).
   - Optional NGO partner code.
4. Connect a thermal printer and barcode scanner over USB. Both are detected automatically.
5. Run **Setup → Test print** to confirm the receipt printer.
6. Run **Setup → Sync now** to seed the offline cache.
7. Hand-off to consular staff.

## Backups

The kiosk syncs its offline cache to the central instance every
60 seconds when connectivity is available. SD-card failure is the
expected loss event — a replacement boots, registers under the same
embassy ID, and resumes within minutes.

## Updates

OTA updates are signed and require a counter-signature from the
trust & safety council ([GOVERNANCE](GOVERNANCE.md)). A failed update
auto-rolls back to the previous image.

## Troubleshooting

| Symptom                  | Likely cause              | Action |
|--------------------------|---------------------------|--------|
| Scanner doesn't read     | USB power-budget low      | Use the powered hub from the BOM |
| Printer prints blank     | Wrong paper width         | Switch to 80 mm thermal paper |
| Verdict spinner forever  | DNS over LTE blocked      | Use 1.1.1.1 / 9.9.9.9 fallbacks |
