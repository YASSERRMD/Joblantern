# ADR-012 — Browser extension architecture

- **Status:** Accepted
- **Date:** 2026-05-26

## Context

Joblantern's web UI is useful, but most users discover suspicious
listings *while browsing job boards*, not after the fact. We need a
surface that meets them where they are.

## Decision

- **Manifest v3 for Chrome / Chromium**, Manifest v2 for Firefox
  (Firefox MV3 is incomplete for service workers in some channels at
  the time of writing).
- **Single shared codebase** under `ext/shared/`, browser-specific
  manifests under `ext/chrome/` and `ext/firefox/`.
- **No bundler, no npm, no framework.** Plain ES2020 + small DOM
  helpers. Mirrors the project-wide "no JS build chain for the
  application" preference (ADR-001).
- **Site registry** pattern: each supported board ships a small
  `sites/<id>.js` file exporting `{ match, extract, anchorSelector }`.
  Adding a board is a one-file change.
- **Background service worker** owns API calls + IndexedDB cache;
  content scripts never talk to the network directly.
- **No `<all_urls>` permission**, no broad `tabs` permission. Only the
  configured job-board hosts.

## Consequences

### Positive

- Reviewable: a Chrome Web Store / AMO reviewer can audit every file
  in minutes.
- No supply-chain surface from npm.
- Per-site extractors are independently fixable when a job board's
  DOM drifts.

### Negative

- DOM selectors are inherently brittle — see Phase 21 follow-up for
  Playwright tests with checked-in fixture HTML.
- MV3 + MV2 forking means two manifests. Mitigated by keeping the
  shared code identical and using `typeof browser !== "undefined" ?
  browser : chrome` to abstract the API.

## References

- WebExtension API: https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions
- Chrome MV3 docs: https://developer.chrome.com/docs/extensions/mv3/
- `ext/PRIVACY.md`, `ext/PERMISSIONS.md` for the store-review docs.
