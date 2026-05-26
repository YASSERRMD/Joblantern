# Joblantern WebExtension — Privacy Notice

## What we send to the Joblantern API

When you load a supported job listing, the extension extracts:

- Listing URL
- Job title and description
- Company name
- Office address as shown on the page
- Recruiter contact details when visible
- The site identifier (e.g. `linkedin`)

This is POSTed to the Joblantern API endpoint **you have configured** in the
extension's Settings page. The default value is `http://localhost:8080`,
which means **nothing leaves your computer** until you point the extension
at a real Joblantern server.

## What we do **not** send

- Your name, email, IP, or any account identifier.
- Cookies from the job board.
- Page content other than the fields listed above.
- Telemetry — unless you opt in via Settings.

## What we store locally

- A 7-day cache of verdicts keyed by listing URL (in IndexedDB).
- Your settings (in `chrome.storage.local`).

You can clear both at any time by removing the extension or opening
`chrome://extensions/` → Joblantern → Storage → Clear.

## What the server stores

That depends on the Joblantern instance you point at. The default
[`SECURITY.md`](https://github.com/yasserrmd/joblantern/blob/main/SECURITY.md)
and [`docs/PRIVACY.md`](https://github.com/yasserrmd/joblantern/blob/main/docs/PRIVACY.md)
in the main repo describe the reference deployment's retention policy
(90 days for verifications, 30 days for audit log).

## Open source

Source: https://github.com/yasserrmd/joblantern (Apache 2.0).
