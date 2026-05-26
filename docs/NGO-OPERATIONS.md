# NGO operations

Tools for NGO caseworkers running Joblantern on behalf of walk-in
clients, hotline callers, and remote case management.

## Kiosk mode

`/kiosk` is a large-font, no-chrome, no-JS form designed for an in-
person clinic tablet. A caseworker turns the device around for the
client; the client pastes the recruiter message and taps **Check**.

The result page at `/kiosk/result/{id}` auto-refreshes every 3 seconds
while the verdict is running, then shows a colour-coded headline and
the top 5 reasons in plain language.

## Paper printout

`/verifications/{id}/print.pdf` renders a single-page A4 PDF with the
risk badge, confidence, top reasons, and the mandatory legal
disclaimer. The PDF is suitable for printing and handing to a client
who does not own a smartphone.

Encoding: Helvetica is built into every PDF reader; we do not embed
fonts (keeps the file ~3 KB).

## What is **not** in v0.1

The Phase 30 roadmap also calls for: IMAP listener for forwarded
emails, weekly digest email, white-label theming, trauma-informed UI
variant, case-management tagging, training mode, hotline call logger.
Those land as follow-up PRs — each is independently shippable on top
of the kiosk + PDF foundation here.

## Locale

The kiosk page is English-only in v0.1. The plumbing to swap copy via
the `i18n` package (Phase 16) is a small additional commit; deferred
to keep this phase narrowly focused.
