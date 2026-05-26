# Joblantern Mobile (Flutter)

```bash
cd mobile
flutter pub get
flutter run -d <device>
```

Initial API base URL is hardcoded in `lib/main.dart` and can be
overridden via `shared_preferences` (a Settings screen lands in a
follow-up commit).

See `../docs/MOBILE.md` for the full scope (offline cache, biometric
lock, deep links, F-Droid metadata, CI builds).
