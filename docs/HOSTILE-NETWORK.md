# Operating Joblantern on hostile networks

A short field guide for NGO operators and individual users in regions
where the network is monitored, throttled, or unreliable.

## Tor onion service

```bash
docker compose -f deploy/tor/docker-compose.yml up -d
cat deploy/tor/data/hidden_service/hostname
```

The hidden service forwards onion port 80 to `host.docker.internal:8080`
where the main `joblantern` binary is listening. Share the resulting
`*.onion` address with your users; they can reach you via Tor Browser
without ever exposing the public-clearnet host.

⚠ This deployment is a **convenience** — it puts the onion service on
the same machine as the Joblantern server. For a real production
target consider running Tor on a dedicated gateway box and using
Whonix-style isolation. See the threat model in `docs/THREAT-MODEL.md`.

## Lite (minimum-bandwidth) mode

`/lite` serves a ~1 KB form: pure HTML, no CSS, no JS, no map tiles,
no service worker. It posts to the same `/verify` endpoint. Useful
on EDGE / 2G connections or via Tor where every byte costs.

## Panic-wipe

`/panic-wipe` (link in the lite UI footer) returns a
`Clear-Site-Data` header plus a small inline script that removes:

- Service-worker registrations
- All caches accessible via the Cache API
- All IndexedDB databases (`joblantern`, `joblantern-pwa`)
- `localStorage` and `sessionStorage`
- Cookies for the site

A user about to walk through a hostile checkpoint can hit the link
and remove every trace of Joblantern from their browser in one tap.

## Reproducible builds (recommended)

For higher-assurance deployments, build the `joblantern` binary from
a pinned `go.sum` and a clean container, and store the SHA256 of the
output. The Dockerfiles in `deploy/` and `cmd/*/` use `-trimpath
-ldflags="-s -w"` so the build is reproducible across hosts.

## Snowflake / domain-fronting

Joblantern itself is a plain HTTP server; any of the standard
censorship-resistant transports (Snowflake, meek, Conjure, etc.) can
sit in front of it. We don't ship a tool for this — the choice is
operator- and jurisdiction-specific. See the Tor Project's pluggable
transports docs.

## See also

- `docs/THREAT-MODEL.md`
- `docs/PRIVACY.md`
- `docs/RETENTION.md`
