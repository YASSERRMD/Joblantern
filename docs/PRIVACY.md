# Privacy

## What Joblantern stores

- The raw submission (URL or pasted text).
- The agent's evidence facts and verdict.
- Optional feedback (thumbs-up/down + comment) on the result page.
- Anonymous identifiers (UUIDs) for verifications and feedback rows.

## What it does not store

- The user's IP address (rate-limit middleware uses hashed values).
- Identity documents — Joblantern never asks for them.
- LLM-side prompt or response logs.

## How long

- Verifications and evidence: 90 days by default.
- Feedback: kept until moderated then either archived (anonymised
  fact joins scam_reports) or deleted.
- Audit log: 30 days, rolled to cold storage.

## Deletion

A signed-in user can delete their submission history through their
account page. NGOs operating their own instance can wipe the tables
at any time via standard SQL.

## Sharing

Joblantern does not sell or share submission data with third parties.
The anonymised scam-feed export only contains: company name, country,
high-level pattern tags. No user content, no recruiter contact details.
