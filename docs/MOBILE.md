# Joblantern mobile (Flutter)

A native iOS / Android client sharing the Joblantern HTTP API.

## What ships in v0.1

- `mobile/` Flutter project with home + result screens.
- `JoblanternApi` HTTP client (verify + poll).
- Material 3 dark theme.

## What's deferred (next mobile branch)

- Settings screen for API endpoint and language.
- Offline-first verdict cache via `drift` (sqlite).
- Biometric lock on history (`local_auth`).
- Deep links from `https://joblantern.example/verifications/<id>` →
  the native screen.
- `flutter_map` + OSM tiles for the result page.
- F-Droid metadata in `mobile/metadata/`.
- GitHub Actions CI build that produces an APK + IPA on tag.

## Why these are deferred

iOS signing and Play Store / F-Droid publication need infra (Apple
Developer account, keystore, F-Droid build server tagging) that lives
outside the agent's reach. The scaffold here is structured so each of
those is a self-contained PR.

## Local dev

```bash
cd mobile
flutter pub get
flutter test
flutter run -d <device>
```
