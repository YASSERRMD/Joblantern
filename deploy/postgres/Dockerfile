# syntax=docker/dockerfile:1.7
#
# Joblantern Postgres image
# -------------------------
# Built FROM the official PostGIS image (which already includes PostgreSQL 16
# and PostGIS 3.4 / 3.5), then adds:
#
#   * pgvector       — embedding storage and HNSW index
#   * pg_trgm        — fuzzy text matching (ships in contrib; we just CREATE it)
#   * citext         — case-insensitive text (ships in contrib; we just CREATE it)
#   * uuid-ossp      — UUID generators (ships in contrib; we just CREATE it)
#
# pg_trgm, citext and uuid-ossp are part of postgresql-contrib which the
# postgis/postgis base image already installs, so no extra apt step is needed
# beyond pgvector. They are enabled by migration 0001.
#
# License: PostgreSQL License (Postgres) + Apache-2.0 (PostGIS) + PostgreSQL-style (pgvector)
#
ARG POSTGRES_VERSION=16
ARG POSTGIS_VERSION=3.4
FROM postgis/postgis:${POSTGRES_VERSION}-${POSTGIS_VERSION}

ARG PGVECTOR_VERSION=v0.8.0

# Install pgvector from source. The postgis/postgis image is Debian-based and
# already has build-essential dependencies for Postgres extensions via apt.
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        build-essential \
        git \
        postgresql-server-dev-all \
        ca-certificates; \
    git clone --depth 1 --branch ${PGVECTOR_VERSION} https://github.com/pgvector/pgvector.git /tmp/pgvector; \
    cd /tmp/pgvector; \
    make OPTFLAGS=""; \
    make install; \
    cd /; rm -rf /tmp/pgvector; \
    apt-get purge -y --auto-remove build-essential git postgresql-server-dev-all; \
    rm -rf /var/lib/apt/lists/*

# Image labels for provenance.
LABEL org.opencontainers.image.title="joblantern-postgres"
LABEL org.opencontainers.image.source="https://github.com/yasserrmd/joblantern"
LABEL org.opencontainers.image.licenses="PostgreSQL AND Apache-2.0"
LABEL org.opencontainers.image.description="PostgreSQL 16 + PostGIS 3.4 + pgvector for Joblantern."
