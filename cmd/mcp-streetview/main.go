// Package main is the joblantern.streetview MCP server.
//
// Tools:
//
//	images_near_point      list images within radius_m of (lat, lon)
//	latest_image_age       age of most recent image near a point
//	_meta_attribution      Mapillary attribution string
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/mapillary"
)

const ServerName = "joblantern.streetview"

const (
	ErrTokenInvalid = "TOKEN_INVALID"
	ErrRateLimited  = "RATE_LIMITED"
	ErrUpstream     = "UPSTREAM_ERROR"
	ErrInvalidArgs  = "INVALID_ARGS"
)

func main() {
	transport := flag.String("transport", getenv("TRANSPORT", "stdio"), "stdio | http")
	addr := flag.String("addr", getenv("ADDR", ":8083"), "HTTP listen addr")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	token := os.Getenv("MAPILLARY_TOKEN")
	if token == "" {
		logger.Warn("MAPILLARY_TOKEN not set; tools will return TOKEN_INVALID")
	}
	cli := mapillary.New(token)

	s := newServer(cli)

	switch *transport {
	case "stdio":
		if err := runStdio(s); err != nil {
			logger.Error("stdio exit", "err", err)
			os.Exit(1)
		}
	case "http":
		if err := runHTTP(s, *addr); err != nil {
			logger.Error("http exit", "err", err)
			os.Exit(1)
		}
	}
}

type imagesArgs struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m,omitempty"`
	Max     int     `json:"max,omitempty"`
}

type imagesResult struct {
	Images []imageOut `json:"images"`
	Code   string     `json:"code,omitempty"`
}

type imageOut struct {
	ID         string `json:"id"`
	CapturedAt string `json:"captured_at,omitempty"`
	ThumbURL   string `json:"thumb_url,omitempty"`
	IsPano     bool   `json:"is_pano,omitempty"`
}

type ageArgs struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m,omitempty"`
}

type ageResult struct {
	LatestCapturedAt string `json:"latest_captured_at,omitempty"`
	AgeDays          int    `json:"age_days"`
	Code             string `json:"code,omitempty"`
}

type attribArgs struct{}
type attribResult struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func newServer(cli *mapillary.Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: ServerName, Version: "0.0.1"}, nil)

	mcp.AddTool(s,
		&mcp.Tool{Name: "images_near_point", Description: "List Mapillary images within radius_m metres of (lat, lon)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a imagesArgs) (*mcp.CallToolResult, imagesResult, error) {
			if cli.Token == "" {
				return errCT(ErrTokenInvalid, "MAPILLARY_TOKEN not set"),
					imagesResult{Code: ErrTokenInvalid}, nil
			}
			imgs, err := cli.ImagesNearPoint(ctx, a.Lat, a.Lon, a.RadiusM, a.Max)
			if err != nil {
				return mapErr(err)
			}
			out := imagesResult{Images: make([]imageOut, 0, len(imgs))}
			for _, im := range imgs {
				out.Images = append(out.Images, imageOut{
					ID:         im.ID,
					CapturedAt: timeStr(im.CapturedTime()),
					ThumbURL:   im.ThumbURL,
					IsPano:     im.IsPano,
				})
			}
			b, _ := json.Marshal(out)
			return okCT(string(b)), out, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "latest_image_age", Description: "Age (in days) of the most recent Mapillary image near (lat, lon)."},
		func(ctx context.Context, _ *mcp.CallToolRequest, a ageArgs) (*mcp.CallToolResult, ageResult, error) {
			if cli.Token == "" {
				return errCT(ErrTokenInvalid, "MAPILLARY_TOKEN not set"),
					ageResult{Code: ErrTokenInvalid}, nil
			}
			imgs, err := cli.ImagesNearPoint(ctx, a.Lat, a.Lon, a.RadiusM, 100)
			if err != nil {
				_, r, _ := mapErrAge(err)
				return errCTAge(err), r, nil
			}
			if len(imgs) == 0 {
				return okCT("no images"), ageResult{AgeDays: -1}, nil
			}
			sort.Slice(imgs, func(i, j int) bool {
				return imgs[i].CapturedAt > imgs[j].CapturedAt
			})
			latest := imgs[0].CapturedTime()
			age := int(time.Since(latest).Hours() / 24)
			return okCT(timeStr(latest)), ageResult{
				LatestCapturedAt: timeStr(latest),
				AgeDays:          age,
			}, nil
		})

	mcp.AddTool(s,
		&mcp.Tool{Name: "_meta_attribution", Description: "Required Mapillary attribution string for any UI displaying these images."},
		func(_ context.Context, _ *mcp.CallToolRequest, _ attribArgs) (*mcp.CallToolResult, attribResult, error) {
			r := attribResult{
				Text: "Imagery © Mapillary (CC BY-SA 4.0)",
				URL:  "https://www.mapillary.com/",
			}
			return okCT(r.Text), r, nil
		})

	return s
}

func mapErr(err error) (*mcp.CallToolResult, imagesResult, error) {
	switch {
	case errors.Is(err, mapillary.ErrTokenInvalid):
		return errCT(ErrTokenInvalid, err.Error()), imagesResult{Code: ErrTokenInvalid}, nil
	case errors.Is(err, mapillary.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), imagesResult{Code: ErrRateLimited}, nil
	default:
		return errCT(ErrUpstream, err.Error()), imagesResult{Code: ErrUpstream}, nil
	}
}

func mapErrAge(err error) (*mcp.CallToolResult, ageResult, error) {
	switch {
	case errors.Is(err, mapillary.ErrTokenInvalid):
		return errCT(ErrTokenInvalid, err.Error()), ageResult{Code: ErrTokenInvalid}, nil
	case errors.Is(err, mapillary.ErrRateLimited):
		return errCT(ErrRateLimited, err.Error()), ageResult{Code: ErrRateLimited}, nil
	default:
		return errCT(ErrUpstream, err.Error()), ageResult{Code: ErrUpstream}, nil
	}
}

func errCTAge(err error) *mcp.CallToolResult { return errCT(ErrUpstream, err.Error()) }

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
		slog.Info("mcp-streetview listening", "addr", addr)
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
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func okCT(t string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: t}}}
}
func errCT(code, msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("[%s] %s", code, msg)}}}
}

func timeStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func getenv(k, fb string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fb
}
