# Security Policy

## Reporting a vulnerability

Please report security vulnerabilities privately to **security@joblantern.org**
(placeholder until the domain is provisioned; in the interim, open a private
GitHub Security Advisory on this repository).

Please include:

- A description of the issue and its impact.
- Steps to reproduce, ideally with a minimal proof-of-concept.
- Affected version (commit SHA or release tag).
- Your name and contact details if you would like to be credited.

## Disclosure window

We aim to acknowledge reports within **72 hours** and to ship a fix or
mitigation within **90 days** of acknowledgement. We will coordinate
public disclosure with you. Please do not file public issues, post to
social media, or share proof-of-concept exploits before the window
closes.

## Scope

In scope:

- The Joblantern Go binaries (`cmd/joblantern`, `cmd/mcp-*`).
- The PostgreSQL schema and migrations.
- The web UI rendered by `cmd/joblantern`.
- The Docker Compose deployment files we publish.

Out of scope:

- Vulnerabilities in upstream dependencies (please report those upstream
  and let us know so we can pin a fixed version).
- Findings that require a malicious operator with shell access to the
  host running Joblantern.
- Rate limiting or volumetric DoS on the public demo instance.

## Safe harbour

Good-faith security research that follows this policy will not result
in legal action from the Joblantern project.
