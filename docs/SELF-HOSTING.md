# Self-hosting Joblantern

A minimal recipe for an NGO or solo operator who wants to run their
own Joblantern instance.

## Requirements

- Linux host with Docker 24+ and Docker Compose v2.
- 2 CPU / 4 GB RAM is enough for the app + Postgres.
- A regional OSM extract from Geofabrik if you want street-level features.
- API tokens (free) for Mapillary, OpenCorporates, OpenRouteService — all
  optional; the agent degrades gracefully when a token is absent.

## Bring-up

```bash
git clone https://github.com/yasserrmd/joblantern.git
cd joblantern
cp .env.example .env
# fill in tokens you have; leave the rest blank
make docker-up
make migrate-up
```

The main app is at `http://localhost:8080`.

## Optional: self-hosted Nominatim + Overpass

```bash
docker compose -f deploy/nominatim/docker-compose.yml up -d
docker compose -f deploy/overpass/docker-compose.yml up -d
```

Set `NOMINATIM_URL` and `OVERPASS_URL` in `.env` to point at these.

## Updating

```bash
git pull
make docker-up
make migrate-up
```

## Backups

Postgres is the single source of truth. Take pg_dumps on whatever
cadence your operation needs:

```bash
docker exec joblantern-postgres pg_dump -U joblantern -Fc joblantern \
    > backups/joblantern-$(date +%F).pgdump
```

## Observability

- `/metrics` — Prometheus format.
- `/healthz` — plain `ok`.
