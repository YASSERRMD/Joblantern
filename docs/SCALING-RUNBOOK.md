# Scaling Runbook

## Target

10M verdicts/day sustained, p99 verdict latency under 8s.

## Footprint (reference)

| Tier            | Instances              | Memory | Notes |
|-----------------|------------------------|--------|-------|
| Web / API       | 12 × 4 vCPU, 8 GB      | 96 GB  | Behind a load balancer + CDN. |
| Workers         | 24 × 8 vCPU, 16 GB     | 384 GB | Queue-driven; scale on queue depth. |
| Postgres primary| 1 × 32 vCPU, 256 GB    | 256 GB | pgbouncer in front. |
| Postgres replicas| 3 × 32 vCPU, 256 GB   | 768 GB | One per region. |
| Cache           | 3 × 4 vCPU, 16 GB      | 48 GB  | DragonflyDB cluster. |
| Object storage  | S3 or MinIO            |        | Evidence artifacts. |

## Checklist for a 10M/day cutover

1. [pgbouncer](../deploy/scale/pgbouncer.ini) deployed.
2. Read replicas online — see [replica router](../internal/scale/replica/router.go).
3. Verifications partitioned by month.
4. Queue cutover from in-process to River.
5. CDN cache rules applied.
6. Load test ([k6](../deploy/scale/loadtest.k6.js)) holds p99 < 8s.
7. Autoscale policies in place.
8. Cost-per-verdict dashboard live.

## Failure modes

- **Queue depth climbs**: scale workers; investigate slow MCP server.
- **p99 spikes only**: cache miss storm; verify hit rate and pre-warm.
- **Disk fills**: rotate to a fresh monthly partition; archive older.

## See also

- [SUSTAINABILITY](SUSTAINABILITY.md)
- [HOSTED-SERVICE](HOSTED-SERVICE.md)
