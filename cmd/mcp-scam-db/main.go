// Package main is the joblantern.scamdb MCP server.
//
// Backed by the scam_reports table (migration 0008). The server is
// the only place where scam-DB writes happen in v1 — the agent reads
// through this server, the moderation tooling writes through it.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerName = "joblantern.scamdb"

const (
	ErrInvalidArgs = "INVALID_ARGS"
	ErrDB          = "DB_ERROR"
	ErrEmbedding   = "EMBEDDING_UNAVAILABLE"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8086"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	dsn := getenv("DATABASE_URL", "postgres://joblantern:joblantern@localhost:5432/joblantern?sslmode=disable")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("open db", "err", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	s := newServer(db)
	switch *transport {
	case "stdio":
		if err := runStdio(s); err != nil {
			logger.Error("stdio", "err", err)
			os.Exit(1)
		}
	case "http":
		if err := runHTTP(s, *addr); err != nil {
			logger.Error("http", "err", err)
			os.Exit(1)
		}
	}
}

type Match struct {
	ID           string    `json:"id"`
	CompanyName  string    `json:"company_name,omitempty"`
	Address      string    `json:"address,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	Domain       string    `json:"domain,omitempty"`
	ReportSource string    `json:"report_source,omitempty"`
	ReportURL    string    `json:"report_url,omitempty"`
	Summary      string    `json:"summary,omitempty"`
	DistanceM    *float64  `json:"distance_m,omitempty"`
	Similarity   *float64  `json:"similarity,omitempty"`
	ReportedAt   time.Time `json:"reported_at,omitempty"`
}

type result struct {
	Matches []Match `json:"matches"`
	Code    string  `json:"code,omitempty"`
}

type phoneArgs struct {
	Phone string `json:"phone"`
	Limit int    `json:"limit,omitempty"`
}
type emailArgs struct {
	Email string `json:"email"`
	Limit int    `json:"limit,omitempty"`
}
type domArgs struct {
	Domain string `json:"domain"`
	Limit  int    `json:"limit,omitempty"`
}
type nameArgs struct {
	Name   string  `json:"name"`
	MinSim float64 `json:"min_sim,omitempty"`
	Limit  int     `json:"limit,omitempty"`
}
type nearArgs struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m,omitempty"`
	Limit   int     `json:"limit,omitempty"`
}
type insertArgs struct {
	CompanyName string  `json:"company_name,omitempty"`
	Address     string  `json:"address,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Email       string  `json:"email,omitempty"`
	Domain      string  `json:"domain,omitempty"`
	Summary     string  `json:"summary,omitempty"`
	Source      string  `json:"source,omitempty"`
	URL         string  `json:"url,omitempty"`
}
type insertResult struct {
	ID   string `json:"id"`
	Code string `json:"code,omitempty"`
}

func newServer(db *sql.DB) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "search_reports_by_phone", Description: "Exact-match scam reports by normalised E.164 phone."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a phoneArgs) (*mcp.CallToolResult, result, error) {
			if a.Phone == "" {
				return errCT(ErrInvalidArgs, "phone required"), result{Code: ErrInvalidArgs}, nil
			}
			rows, err := query(ctx, db, `
				SELECT id::text, company_name, address, phone, email, domain, report_source, report_url, summary, reported_at
				  FROM scam_reports WHERE phone = $1
				 ORDER BY reported_at DESC NULLS LAST LIMIT $2`, a.Phone, lim(a.Limit, 20))
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			r := result{Matches: rows}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "search_reports_by_email", Description: "Exact-match scam reports by email (or email-domain prefix '%@example')."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a emailArgs) (*mcp.CallToolResult, result, error) {
			if a.Email == "" {
				return errCT(ErrInvalidArgs, "email required"), result{Code: ErrInvalidArgs}, nil
			}
			rows, err := query(ctx, db, `
				SELECT id::text, company_name, address, phone, email, domain, report_source, report_url, summary, reported_at
				  FROM scam_reports WHERE email = $1 OR email ILIKE '%@' || $1
				 ORDER BY reported_at DESC NULLS LAST LIMIT $2`, a.Email, lim(a.Limit, 20))
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			r := result{Matches: rows}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "search_reports_by_domain", Description: "Scam reports for an exact domain or any subdomain."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a domArgs) (*mcp.CallToolResult, result, error) {
			if a.Domain == "" {
				return errCT(ErrInvalidArgs, "domain required"), result{Code: ErrInvalidArgs}, nil
			}
			rows, err := query(ctx, db, `
				SELECT id::text, company_name, address, phone, email, domain, report_source, report_url, summary, reported_at
				  FROM scam_reports WHERE domain = $1 OR domain ILIKE '%.' || $1
				 ORDER BY reported_at DESC NULLS LAST LIMIT $2`, a.Domain, lim(a.Limit, 20))
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			r := result{Matches: rows}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "search_reports_by_company_name", Description: "Trigram fuzzy company-name match. min_sim defaults to 0.4."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a nameArgs) (*mcp.CallToolResult, result, error) {
			if a.Name == "" {
				return errCT(ErrInvalidArgs, "name required"), result{Code: ErrInvalidArgs}, nil
			}
			minSim := a.MinSim
			if minSim <= 0 || minSim > 1 {
				minSim = 0.4
			}
			rs, err := db.QueryContext(ctx, `
				SELECT id::text, company_name, address, phone, email, domain, report_source, report_url, summary, reported_at,
				       similarity(company_name, $1) AS sim
				  FROM scam_reports
				 WHERE company_name % $1 AND similarity(company_name, $1) >= $2
				 ORDER BY sim DESC
				 LIMIT $3`, a.Name, minSim, lim(a.Limit, 20))
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			defer func() { _ = rs.Close() }()
			matches, err := scanWithSim(rs)
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			r := result{Matches: matches}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "search_reports_near_point", Description: "Scam reports within radius_m metres of (lat, lon)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a nearArgs) (*mcp.CallToolResult, result, error) {
			radius := a.RadiusM
			if radius <= 0 {
				radius = 500
			}
			rs, err := db.QueryContext(ctx, `
				SELECT id::text, company_name, address, phone, email, domain, report_source, report_url, summary, reported_at,
				       ST_Distance(address_geom, ST_SetSRID(ST_MakePoint($1::float8, $2::float8), 4326)::geography) AS dist
				  FROM scam_reports
				 WHERE address_geom IS NOT NULL
				   AND ST_DWithin(address_geom, ST_SetSRID(ST_MakePoint($1::float8, $2::float8), 4326)::geography, $3::float8)
				 ORDER BY dist ASC
				 LIMIT $4`, a.Lon, a.Lat, float64(radius), lim(a.Limit, 20))
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			defer func() { _ = rs.Close() }()
			matches, err := scanWithDist(rs)
			if err != nil {
				return errCT(ErrDB, err.Error()), result{Code: ErrDB}, nil
			}
			r := result{Matches: matches}
			return okCT(jstr(r)), r, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "insert_report", Description: "Internal-only insert; the moderation pipeline calls this after a feedback row is approved."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a insertArgs) (*mcp.CallToolResult, insertResult, error) {
			var id string
			var lat, lon any
			if a.Lat != 0 || a.Lon != 0 {
				lat, lon = a.Lat, a.Lon
			}
			err := db.QueryRowContext(ctx, `
				INSERT INTO scam_reports
					(company_name, address, address_geom, phone, email, domain, summary, report_source, report_url)
				VALUES ($1, $2,
					CASE WHEN $3::float8 IS NULL OR $4::float8 IS NULL THEN NULL
					     ELSE ST_SetSRID(ST_MakePoint($4::float8, $3::float8), 4326)::geography END,
					$5, $6, $7, $8, $9, $10)
				RETURNING id::text`,
				a.CompanyName, a.Address, lat, lon,
				a.Phone, a.Email, a.Domain, a.Summary, a.Source, a.URL,
			).Scan(&id)
			if err != nil {
				return errCT(ErrDB, err.Error()), insertResult{Code: ErrDB}, nil
			}
			return okCT(id), insertResult{ID: id}, nil
		})

	return s
}

// --- helpers ---

func query(ctx context.Context, db *sql.DB, q string, args ...any) ([]Match, error) {
	rs, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rs.Close() }()
	var out []Match
	for rs.Next() {
		var m Match
		var company, address, phone, email, domain, src, url, summary sql.NullString
		var reportedAt sql.NullTime
		if err := rs.Scan(&m.ID, &company, &address, &phone, &email, &domain, &src, &url, &summary, &reportedAt); err != nil {
			return nil, err
		}
		assign(&m, company, address, phone, email, domain, src, url, summary)
		if reportedAt.Valid {
			m.ReportedAt = reportedAt.Time
		}
		out = append(out, m)
	}
	return out, rs.Err()
}

func scanWithSim(rs *sql.Rows) ([]Match, error) {
	var out []Match
	for rs.Next() {
		var m Match
		var company, address, phone, email, domain, src, url, summary sql.NullString
		var reportedAt sql.NullTime
		var sim float64
		if err := rs.Scan(&m.ID, &company, &address, &phone, &email, &domain, &src, &url, &summary, &reportedAt, &sim); err != nil {
			return nil, err
		}
		assign(&m, company, address, phone, email, domain, src, url, summary)
		if reportedAt.Valid {
			m.ReportedAt = reportedAt.Time
		}
		s := sim
		m.Similarity = &s
		out = append(out, m)
	}
	return out, rs.Err()
}

func scanWithDist(rs *sql.Rows) ([]Match, error) {
	var out []Match
	for rs.Next() {
		var m Match
		var company, address, phone, email, domain, src, url, summary sql.NullString
		var reportedAt sql.NullTime
		var dist float64
		if err := rs.Scan(&m.ID, &company, &address, &phone, &email, &domain, &src, &url, &summary, &reportedAt, &dist); err != nil {
			return nil, err
		}
		assign(&m, company, address, phone, email, domain, src, url, summary)
		if reportedAt.Valid {
			m.ReportedAt = reportedAt.Time
		}
		d := dist
		m.DistanceM = &d
		out = append(out, m)
	}
	return out, rs.Err()
}

func assign(m *Match, company, address, phone, email, domain, src, url, summary sql.NullString) {
	if company.Valid {
		m.CompanyName = company.String
	}
	if address.Valid {
		m.Address = address.String
	}
	if phone.Valid {
		m.Phone = phone.String
	}
	if email.Valid {
		m.Email = email.String
	}
	if domain.Valid {
		m.Domain = domain.String
	}
	if src.Valid {
		m.ReportSource = src.String
	}
	if url.Valid {
		m.ReportURL = url.String
	}
	if summary.Valid {
		m.Summary = summary.String
	}
}

func runStdio(s *mcp.Server) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	session, err := s.Connect(ctx, &mcp.StdioTransport{}, nil)
	if err != nil {
		return err
	}
	session.Wait()
	return nil
}

func runHTTP(s *mcp.Server, addr string) error {
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return s }, nil)
	srv := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("mcp-scam-db listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(sctx)
}

func okCT(t string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: t}}}
}
func errCT(code, msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("[%s] %s", code, msg)}}}
}
func jstr(v any) string { b, _ := json.Marshal(v); return string(b) }
func lim(n, d int) int {
	if n <= 0 {
		return d
	}
	return n
}
func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
