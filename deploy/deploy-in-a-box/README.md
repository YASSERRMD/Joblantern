# Deploy-In-A-Box

Pre-configured images so an NGO partner can spin up an instance
without specialised ops.

## Images

| Cloud           | Image / snapshot                                    |
|-----------------|------------------------------------------------------|
| Hetzner Cloud   | `joblantern-${version}-x86`                          |
| Digital Ocean   | snapshot `joblantern-${version}`                     |
| AWS (community) | AMI `ami-joblantern-${version}` (eu-central-1 first) |

## Boot

```bash
# After provisioning a $4 droplet / VM:
ssh root@<ip>
joblantern-setup --slug my-ngo --region eu --email ops@my-ngo.org
```

The setup script generates ed25519 keys, requests a TLS cert via
ACME, and seeds the database with the standard rule pack.

## What's pre-installed

- Joblantern primary binary + MCP servers
- pgbouncer + Postgres
- Caddy with ACME + HSTS preload
- ufw with the published deny-by-default ruleset

## What you still need to do

1. Point a DNS A record at the VM IP.
2. Re-run `joblantern-setup` once DNS resolves so ACME succeeds.
3. Walk through the [first-verdict](../../internal/capacity/curriculum/modules.go) tutorial.

## Sizing

For < 5,000 verdicts/month a single 4 GB VM is enough. Scale up by
following [SCALING-RUNBOOK](../../docs/SCALING-RUNBOOK.md).
