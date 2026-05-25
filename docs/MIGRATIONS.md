# Migrations workflow

Joblantern uses [`goose`](https://github.com/pressly/goose) for SQL
migrations, wrapped in a small project-owned CLI under
`cmd/goose-migrate` that imports goose as a library and pulls in only
the `pgx/stdlib` Postgres driver. This keeps our dependency surface
narrow and free of database drivers we will never ship.

## Make targets

| Target | Effect |
|---|---|
| `make migrate-up` | Apply all pending migrations against `$DATABASE_URL`. |
| `make migrate-down` | Roll back the most recent migration. |
| `make migrate-status` | Show every migration and whether it has been applied. |
| `make migrate-reset` | Roll all migrations back to zero. **Development only.** |
| `NAME=add_thing make migrate-create` | Scaffold a new SQL migration file in `migrations/`. |

All targets read `$DATABASE_URL` from the environment (defaulted to the
local Compose Postgres in the Makefile).

## File layout

```
migrations/
  0001_extensions.sql
  0002_users.sql
  0003_sessions.sql
  ...
```

Goose uses the leading numeric prefix as the version. Migrations are
applied in numeric order. Names use snake_case.

## Authoring rules

1. **One migration per atomic change.** A migration that creates a
   table and adds an index is fine; a migration that creates four
   unrelated tables is not.
2. **Always include a `-- +goose Down` section** that fully reverses
   the `Up`. Joblantern is open-source — operators must be able to
   roll back safely.
3. **Wrap each migration in an implicit transaction** (the default).
   For statements that cannot run inside a transaction (e.g.
   `CREATE INDEX CONCURRENTLY`, `CREATE TYPE ... AS ENUM` followed by
   `ALTER TYPE ... ADD VALUE` in the same migration), add
   `-- +goose NO TRANSACTION` at the top and document why.
4. **Never edit a migration after it has shipped to `main`.** Add a
   new migration instead. The append-only convention is what makes
   migrations replayable on a fresh database.
5. **Test the rollback.** `make migrate-down` after every new
   migration during development.

## Local quickstart

```bash
make docker-up           # boot the joblantern postgres container
make migrate-up          # apply all migrations
make migrate-status
```

## Continuous integration

CI runs migrations against an ephemeral PostGIS+pgvector container via
`testcontainers-go` (see Phase 02 task 21). A green CI run proves both
Up and reproducibility from a blank database.

## Why a wrapper instead of the `goose` CLI binary?

`go install github.com/pressly/goose/v3/cmd/goose@latest` pulls in
drivers for MSSQL, MySQL, ClickHouse, Vertica, YDB, Turso, libsql, and
others — none of which we ship. Importing the goose library directly
with only `pgx/stdlib` keeps the dependency tree small, faster to
build, and easier to license-audit.
