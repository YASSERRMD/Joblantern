# CDN Recipe

Static assets, OpenStreetMap-style map tiles, and the GraphQL
explorer are served through a CDN. Any of Cloudflare, Fastly, or
Bunny works; the recipe below uses Cloudflare because the cache
rules are easy to express in `_headers`.

## Paths

| Path                          | Cache TTL | Notes |
|-------------------------------|-----------|-------|
| `/static/*`                   | 7 days    | Versioned by content hash. |
| `/tiles/*`                    | 30 days   | Tile sources tested for staleness via ETag. |
| `/api/v1/openapi.json`        | 5 min     | New ones ship per release. |
| `/graphql/explorer`           | 1 day     | Plain HTML, hand-edit triggers a purge. |

## Headers

```
/static/*
  Cache-Control: public, max-age=604800, immutable

/tiles/*
  Cache-Control: public, max-age=2592000, stale-while-revalidate=86400
```

## Purge

`make cdn-purge` runs after a release.
