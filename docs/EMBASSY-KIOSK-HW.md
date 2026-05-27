# Embassy Kiosk — Hardware Spec

The kiosk targets a low-cost, ruggedised, low-bandwidth installation
that a consular IT team can deploy in 30 minutes.

## Minimum bill of materials

| Component        | Spec                                              |
|------------------|---------------------------------------------------|
| SBC / Mini PC    | Quad-core ARMv8 (e.g. RPi 5 / NanoPi R6S), 4 GB RAM |
| Touchscreen      | 10–15", 1280×800, capacitive, IP54                 |
| Receipt printer  | ESC/POS thermal, 80 mm, USB                        |
| QR/barcode scanner | USB HID, 1D + 2D                                 |
| Storage          | 64 GB eMMC + 32 GB SD as offline cache             |
| Network          | Ethernet preferred, LTE failover via USB modem     |
| UPS              | 30-min runtime mini-UPS                            |
| Enclosure        | Steel kiosk with anti-tamper screws                |

## Optional

- Microphone + speaker for the Phase 33 voice mode.
- Stripe of LED status indicators visible from a queue.

## Why this matters

Many origin-country embassies sit in older buildings with patchy
power and bandwidth. A 30-minute install path means a partner NGO
can volunteer to deploy one without specialised IT staff.
