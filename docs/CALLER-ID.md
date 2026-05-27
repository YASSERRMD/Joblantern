# Real-Time Caller ID

Many scams happen by phone. A user gets a call from a "recruiter" and
must decide in 60 seconds. Phase 42 adds a Truecaller-style overlay
backed by the Joblantern scam database.

## Android

[JoblanternCallScreeningService](../mobile/android/callscreening/JoblanternCallScreeningService.kt)
extends Android's `CallScreeningService`. On an incoming call it:

1. Hashes the phone digits locally using a region-scoped pepper.
2. Looks up the hash against the on-device cache (sub-second).
3. Renders a translucent green / yellow / red banner via [Overlay](../mobile/android/callscreening/Overlay.kt).
4. If the network is reachable and the device is not in battery saver, also pulls a fresh server-side verdict in parallel.

## iOS

[JoblanternCallDirectoryHandler](../mobile/ios/callkit/JoblanternCallDirectoryHandler.swift)
ships a Call Directory Extension. Apple does not permit a full
pre-pickup overlay, so we ship best-effort: caller label and a
block-list for high-confidence scams.

## Privacy

- The raw number never leaves the device.
- The HMAC pepper is region-scoped — devices in the same region
  generate the same hash and can correlate; a database leak cannot
  reverse to numbers.

## Performance

Target: red overlay visible within 1 second of the call ringing.
The hash → cache → overlay path runs in tens of milliseconds. The
network round-trip is a parallel refresh, not a blocking step.

## See also

- [PRIVACY](PRIVACY.md)
- [MOBILE](MOBILE.md)
