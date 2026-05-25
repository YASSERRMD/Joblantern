# Self-hosted Overpass

`mcp-address` uses Overpass to classify the land-use / building type
around a coordinate (residential vs commercial vs mixed).

```bash
docker compose -f deploy/overpass/docker-compose.yml up -d
```

Default region: Monaco. Override `OVERPASS_PLANET_URL` for production.

Connect via `OVERPASS_URL=http://localhost:8089/api/interpreter`.

## License

OSM data — ODbL. Display attribution wherever derived results are
shown to users.
