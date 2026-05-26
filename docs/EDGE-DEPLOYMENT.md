# Edge deployment (Raspberry Pi / mini-PC)

Joblantern runs on small offline boxes. Useful for NGO clinics in
low-connectivity regions, or for users who want every verification
verdict computed and explained locally.

## Hardware

| Box | Notes |
|---|---|
| Raspberry Pi 5 / 8 GB | Works with Qwen2.5-3B (Q4_K_M). Verdicts return in 30-60s. |
| Mini-PC, 16 GB | Phi-3-mini or Llama-3.2-3B comfortable. <20s per verdict. |
| Old laptop, 8 GB | Acceptable for caseworker desk use. |

ARM64 and AMD64 are both supported (the release workflow ships both
arches to ghcr.io).

## Bring-up

```bash
# 1) Bring the edge stack up.
docker compose -f deploy/edge/docker-compose.yml up -d

# 2) Pull the default model (~2 GB).
docker exec -it joblantern-edge-ollama-1 \
    ollama pull qwen2.5:3b

# 3) Apply migrations.
make migrate-up DATABASE_URL="postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable"

# 4) Visit http://localhost:8080
```

## Offline mode

`JOBLANTERN_OFFLINE=1` (set in the edge compose) instructs the agent
to skip MCP servers that require external network access (street view,
domain WHOIS / Wayback, OpenCorporates, OpenRouteService). The
verdict will be lower-confidence but still grounded in:

- `joblantern.pattern` (local rules pack)
- `joblantern.law` (bundled jurisdiction rules)
- `joblantern.salary` (bundled bands)
- `joblantern.scamdb` (local Postgres)
- The LLM running on the box

## Model recommendations

| Model | Size (Q4_K_M) | Notes |
|---|---|---|
| `qwen2.5:3b` | 1.9 GB | Default. Multilingual. Apache-2.0. |
| `phi3:mini` | 2.3 GB | Strong reasoning, English-leaning. MIT. |
| `llama3.2:3b` | 2.0 GB | Apache-2.0 derivative; check Meta's terms. |
| `nomic-embed-text` | 274 MB | Sentence embeddings for `mcp-scam-db` semantic search. |

Quantisation guidance: Q4_K_M is the cost/quality sweet spot. Q5_K_M
adds ~25% memory for marginal quality gain. Q8 is not worth it on
edge.

## Power-aware behaviour (planned)

A follow-up PR will read battery state on Linux laptops and skip
heavy MCP calls when on battery + low. Not in v0.1.

## Caveats

- The default Ollama install accepts connections on `0.0.0.0:11434`
  inside the compose network. Production edge deployments should run
  the bundle inside a locked-down network namespace or only bind
  loopback.
- Model files live in the `ollama` Docker volume. Back them up with
  the rest of the stack.
