// Package main is a minimal goose CLI wrapper that supports only the Postgres
// dialect via pgx/stdlib. Using goose as a library (rather than `go install
// github.com/pressly/goose/v3/cmd/goose`) keeps our dependency graph small,
// permissive, and free of database drivers we will never use.
//
// Usage:
//
//	go run ./cmd/goose-migrate up
//	go run ./cmd/goose-migrate down
//	go run ./cmd/goose-migrate status
//	go run ./cmd/goose-migrate create <name> [sql|go]
//
// Connection string and migrations dir come from environment:
//
//	DATABASE_URL        (default: postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable)
//	MIGRATIONS_DIR      (default: migrations)
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	defaultDSN = "postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable"
	defaultDir = "migrations"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "goose-migrate:", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("goose-migrate", flag.ContinueOnError)
	dsn := fs.String("dsn", getenv("DATABASE_URL", defaultDSN), "Postgres connection string")
	dir := fs.String("dir", getenv("MIGRATIONS_DIR", defaultDir), "migrations directory")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	args := fs.Args()
	if len(args) == 0 {
		return errors.New("usage: goose-migrate [flags] <command> [args...]")
	}
	command := args[0]
	commandArgs := args[1:]

	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() { _ = db.Close() }()

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	if err := goose.RunContext(context.Background(), command, db, *dir, commandArgs...); err != nil {
		return fmt.Errorf("goose %s: %w", command, err)
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
