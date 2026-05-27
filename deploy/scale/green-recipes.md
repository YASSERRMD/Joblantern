# Green Deployment Recipes

## Cloudflare Workers

The verdict-presentation surface (the read-only `/v/<id>` page) is a
Worker. It scales to zero between requests and is colocated with the
user. Carbon cost per request is approximately the smallest of any
deployment option we support.

```toml
# wrangler.toml (excerpt)
name = "joblantern-verdict-view"
main = "src/worker.ts"
compatibility_date = "2026-01-01"

[triggers]
crons = []   # no batch on Workers
```

## Fly Machines (scale-to-zero)

The agent workers run on Fly Machines with `auto_stop_machines = true`
so idle regions sleep. Cold-start under one second is acceptable
for batch jobs; verdict requests use a warm pool.

```toml
# fly.toml (excerpt)
[machines]
auto_stop_machines = true
auto_start_machines = true
min_machines_running = 0
```

## Self-host

Use the standard `deploy/compose/*.yml` files. Set
`JL_GREEN_SCHEDULER=on` to defer batch jobs to low-carbon hours.
