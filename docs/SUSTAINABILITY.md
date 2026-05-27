# Sustainability

Operate the service in a way that respects the carbon cost of every
verdict — without compromising the user-facing latency we promise.

## Levers

1. **Eco-scheduling** — batch jobs (red-team, training, exports) defer to low-carbon hours via [scheduler](../internal/greenops/scheduler/eco.go).
2. **Caching** — see [strategy](../internal/greenops/cache/strategy.go).
3. **Compact models** — small LLMs by default, bigger ones only as tie-breakers.
4. **Colocation choice** — [greenest region](../internal/greenops/colocation/region.go) within data-residency constraints.
5. **Scale-to-zero** — Cloudflare Workers + Fly Machines where applicable.

## Footprint accounting

[footprint](../internal/greenops/footprint/report.go) emits per-verdict grams of CO2-eq. The public [sustainability page](../internal/greenops/page/public.go) renders the monthly total.

## Green SLA

User-facing paths (`submit`, `view`) keep the 8s p99 target. Eco-deferrable paths (`aggregate`, batch jobs) accept relaxed latency in exchange for being scheduled into low-carbon windows.

## See also

- [SCALING-RUNBOOK](SCALING-RUNBOOK.md)
- [TRANSPARENCY](TRANSPARENCY.md)
