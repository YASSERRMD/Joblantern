# Joblantern WebExtension

A Manifest v3 / v2 browser extension that adds an inline scam-risk verdict
badge to job listings on supported boards.

## Supported sites

- LinkedIn (`linkedin.com/jobs/view/*`, `linkedin.com/jobs/collections/*`)
- Indeed (any `*.indeed.com/viewjob*`)
- Bayt (`bayt.com/{en,ar}/...`)
- Naukrigulf (`naukrigulf.com/job-listing-*`)
- GulfTalent (`gulftalent.com/{cc}/jobs/*`)
- Jobstreet (any `*.jobstreet.com/job/*`)

Adding a new site is one file: `ext/shared/sites/<name>.js` exporting a
`{ id, name, match, extract, anchorSelector }` to `window.Joblantern`.
The MV3 + MV2 manifests must also be updated to add the new host pattern.

## Install (developer mode)

### Chrome / Chromium

1. Build: `make ext-package` (produces `ext/dist/joblantern-chrome-0.1.0.zip`).
2. `chrome://extensions/` → enable **Developer mode**.
3. Drag the zip onto the page (or click *Load unpacked* and select the
   unzipped staging directory).

### Firefox

1. Build: `make ext-package`.
2. `about:debugging#/runtime/this-firefox` → *Load Temporary Add-on…* →
   select `ext/firefox/manifest.json` (or load the unzipped staging dir).

For permanent installation, the same zip must be uploaded to AMO for
review.

## Configuration

Open the extension's *Settings* page (right-click the icon → Options).

| Setting | Default | Notes |
|---|---|---|
| API endpoint | `http://localhost:8080` | Your Joblantern instance |
| API key | empty | Optional; sent as `X-Joblantern-API-Key` |
| Cache verdicts for (days) | `7` | IndexedDB per-URL cache |
| Telemetry | off | Opt-in only |

## How it works

1. A content script runs on supported listing pages.
2. It picks the right extractor from the site registry and pulls a
   minimal submission (URL, title, company, address, recruiter contact,
   text).
3. The extractor calls `chrome.runtime.sendMessage` to the background
   service worker, which:
   - checks IndexedDB for a fresh cached verdict;
   - otherwise `POST`s `/api/v1/verify`, then polls
     `/api/v1/verifications/{id}` until completion.
4. The verdict is rendered as a coloured pill next to the job title.
   Hovering shows the top reasons; clicking opens the full report at the
   configured API base.

## Tests

The extension targets fixed CSS selectors per site. Selectors break;
add a fixture HTML file under `ext/test/fixtures/<site>.html` and a
Playwright test (Phase 21 PR opens the slot but does not run Playwright
against live sites — those tests fall under follow-up work because they
need browser binaries and per-site review).

## License & privacy

- License: Apache 2.0.
- See `ext/PRIVACY.md` and `ext/PERMISSIONS.md` for the store-review docs.
