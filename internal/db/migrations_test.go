//go:build integration
// +build integration

// Package joblanterndb integration tests.
//
// These tests are guarded by the `integration` build tag because they
// require Docker (or a CI runner with Docker) to spin up Postgres.
//
// Run with:
//
//	make test-integration
//	# or
//	go test -tags=integration ./internal/db/...
//
// The test uses a prebuilt PostGIS + pgvector image so it runs in
// seconds rather than recompiling pgvector each time. The project's own
// deploy/postgres/Dockerfile is validated separately by `docker compose
// build` (Phase 02 verification + CI in Phase 18).
package joblanterndb_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const testImage = "imresamu/postgis-pgvector:16-3.4"

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find go.mod walking up from %s", wd)
		}
		dir = parent
	}
}

func TestMigrations_Up(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	container, err := tcpostgres.Run(ctx,
		testImage,
		tcpostgres.WithDatabase("joblantern"),
		tcpostgres.WithUsername("joblantern"),
		tcpostgres.WithPassword("joblantern"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set dialect: %v", err)
	}
	migrationsDir := filepath.Join(repoRoot(t), "migrations")
	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("goose up: %v", err)
	}

	// Every table from migrations 0002–0010 should exist.
	wantTables := []string{
		"users",
		"sessions",
		"verifications",
		"evidence_facts",
		"scam_reports",
		"jurisdictions",
		"mcp_audit_log",
	}
	for _, table := range wantTables {
		var n int
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			  FROM information_schema.tables
			 WHERE table_schema='public' AND table_name=$1`, table).Scan(&n)
		if err != nil {
			t.Fatalf("inspect table %s: %v", table, err)
		}
		if n != 1 {
			t.Errorf("expected table %s to exist (got count=%d)", table, n)
		}
	}

	// PostGIS, pgvector and pg_trgm should all be live.
	wantExts := []string{"postgis", "vector", "pg_trgm", "uuid-ossp", "citext"}
	for _, ext := range wantExts {
		var present bool
		if err := db.QueryRowContext(ctx,
			`SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname=$1)`,
			ext,
		).Scan(&present); err != nil {
			t.Fatalf("check extension %s: %v", ext, err)
		}
		if !present {
			t.Errorf("extension %s missing", ext)
		}
	}

	// Spatial + vector indexes must have been created on scam_reports.
	wantIndexes := []string{
		"scam_reports_address_geom_idx",
		"scam_reports_company_name_trgm_idx",
		"scam_reports_embedding_hnsw_idx",
	}
	for _, idx := range wantIndexes {
		var present bool
		if err := db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1
				  FROM pg_indexes
				 WHERE schemaname='public' AND indexname=$1
			)`, idx).Scan(&present); err != nil {
			t.Fatalf("check index %s: %v", idx, err)
		}
		if !present {
			t.Errorf("index %s missing", idx)
		}
	}

	// Round-trip a row through scam_reports to prove the column types
	// behave: PostGIS geography input, pgvector embedding input.
	if _, err := db.ExecContext(ctx, `
		INSERT INTO scam_reports
			(company_name, address, address_geom, phone, email, domain, embedding)
		VALUES
			($1, $2,
			 ST_SetSRID(ST_MakePoint($3::float8, $4::float8), 4326)::geography,
			 $5, $6, $7,
			 $8::vector)`,
		"Fake Co", "1 Test St",
		55.2708, 25.2048, // Dubai-ish
		"+971500000000", "scam@fake.example", "fake.example",
		vectorLiteral(384),
	); err != nil {
		t.Fatalf("insert scam_report: %v", err)
	}

	// Spatial neighbour query should find the row we just inserted.
	var hits int
	if err := db.QueryRowContext(ctx, `
		SELECT COUNT(*)::int
		  FROM scam_reports
		 WHERE ST_DWithin(
			   address_geom,
			   ST_SetSRID(ST_MakePoint($1::float8, $2::float8), 4326)::geography,
			   $3::float8)`,
		55.2708, 25.2048, 1000.0,
	).Scan(&hits); err != nil {
		t.Fatalf("spatial query: %v", err)
	}
	if hits != 1 {
		t.Errorf("expected 1 nearby scam_report, got %d", hits)
	}

	// Rolling everything back must succeed.
	if err := goose.Reset(db, migrationsDir); err != nil {
		t.Fatalf("goose reset: %v", err)
	}
}

// vectorLiteral returns a pgvector literal of the given dimensionality
// filled with zeroes — enough for type-system smoke testing without
// pulling pgvector-go just for one test value.
func vectorLiteral(dim int) string {
	b := make([]byte, 0, 2+dim*2)
	b = append(b, '[')
	for i := 0; i < dim; i++ {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, '0')
	}
	b = append(b, ']')
	return string(b)
}
