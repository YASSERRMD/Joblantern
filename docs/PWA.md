# Progressive Web App

Joblantern is installable on mobile and desktop browsers as a PWA. The
service worker handles offline reads of cached verdicts and queues
submissions for replay when the user reconnects.

## What's installable

- Add-to-home-screen on iOS Safari / Android Chrome / desktop Chromium.
- Standalone display mode.
- Theme colour `#11141c` (matches the dark UI).

## Caching strategy

| Path | Strategy |
|---|---|
| `/static/*` | Cache-first |
| `/` | Stale-while-revalidate |
| `/verifications/{id}` | Network-first, falls back to cache or `/static/offline.html` |
| Everything else | Pass through (always network) |

## Background sync

Failed `POST /api/v1/verify` submissions are stored in IndexedDB
(`joblantern-pwa.queue`) and retried by the service worker when the
browser fires the `joblantern-verify-queue` sync event. The current
templ form does not yet queue on POST failure — that wiring is small
and lands as a v0.2 follow-up alongside an offline-first submission UI.

## Web push

The `push` event handler is implemented; VAPID configuration on the
server side (subscription endpoint, key generation, push send) is
deferred to a follow-up commit on this branch family (the `web-push`
Go library + VAPID key generation is non-trivial and warrants its
own PR).

## Lighthouse

Target: ≥ 90 on Performance and PWA. Add `make lighthouse` once the
extension PR's Playwright infra is in (browser binaries needed).

## Icons

Placeholder icon stubs ship under `internal/web/static/icons/`. Real
icons are commissioned art; v0.1 ships transparent PNGs so the
manifest validates.
