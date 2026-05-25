# Self-hosted Nominatim

Joblantern's `mcp-address` server uses Nominatim for geocoding and
reverse-geocoding. The public Nominatim API has a strict acceptable-
use policy (max 1 req/sec, no resale, no heavy use). **Production
deployments MUST self-host.**

## Quickstart

```bash
docker compose -f deploy/nominatim/docker-compose.yml up -d
```

Default region: Monaco (`monaco-latest.osm.pbf`, ~1 MB, ready in a few
minutes). Override with:

```bash
PBF_URL=https://download.geofabrik.de/asia/india-latest.osm.pbf \
REPLICATION_URL=https://download.geofabrik.de/asia/india-updates/ \
docker compose -f deploy/nominatim/docker-compose.yml up -d
```

## Region size table (approximate, cold-start)

| Region | PBF size | Import time | Disk |
|---|---|---|---|
| Monaco | 1 MB | ~3 min | <1 GB |
| UAE | 70 MB | ~30 min | 8 GB |
| India | 850 MB | ~3 h | 60 GB |
| Asia | 12 GB | ~24 h | 600 GB |
| Planet | 80 GB | days | several TB |

## Connecting the address MCP

Set in your environment / `.env`:

```
NOMINATIM_URL=http://localhost:8088
```

## License

OpenStreetMap data is ODbL. Attribution required wherever address
results are displayed.
