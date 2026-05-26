# Joblantern WebExtension — Permissions Justification

The Joblantern extension requests the following permissions. Each is
strictly minimised. This document is intended for Chrome Web Store and
addons.mozilla.org reviewers.

## `storage`

We use `chrome.storage.local` to persist:
- User-chosen API endpoint and (optional) API key.
- Cache TTL preference.

No identifying information is stored under this permission.

## `activeTab`

We do not use broad `tabs` permissions. `activeTab` is only used by the
popup to read the current tab's URL so we can tell the user whether the
extension is active on the open page.

## Host permissions

Limited to the listing pages of supported job boards:

- `https://www.linkedin.com/*` (job pages only via content_scripts matches)
- `https://*.indeed.com/*` (only `viewjob*`)
- `https://www.bayt.com/*`
- `https://www.naukrigulf.com/*`
- `https://www.gulftalent.com/*`
- `https://*.jobstreet.com/*`

We **do not** request `<all_urls>`. We never read pages outside these
hosts. We never inject scripts into bank, payment, or government sites.

## Network requests

The extension only contacts the user-configured Joblantern API endpoint
(default `http://localhost:8080`). It does not call any third-party
servers, analytics, or ad networks.

## Telemetry

Disabled by default. The Settings page exposes a single opt-in checkbox.
When enabled, the extension reports anonymous counters (verifications
attempted, verifications completed) to the same configured API endpoint.

## Source

Full source is at https://github.com/yasserrmd/joblantern under the
Apache 2.0 license. Reviewers can verify the claims above against the
content scripts under `ext/shared/`.
