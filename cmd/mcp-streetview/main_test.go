package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/mapillary"
)

func TestImagesNearPoint_OK(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[{"id":"1","captured_at":1700000000000,"thumb_1024_url":"https://x/y.jpg"}]}`))
	}))
	defer upstream.Close()

	cli := mapillary.New("t")
	cli.BaseURL = upstream.URL
	srv := newServer(cli)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, err := srv.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ss.Close() }()
	c := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, err := c.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "images_near_point",
		Arguments: map[string]any{"lat": 25.2, "lon": 55.27, "radius_m": 50, "max": 5},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}

func TestImagesNearPoint_NoToken(t *testing.T) {
	cli := mapillary.New("")
	srv := newServer(cli)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, _ := srv.Connect(ctx, st, nil)
	defer func() { _ = ss.Close() }()
	c := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, _ := c.Connect(ctx, ct, nil)
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "images_near_point",
		Arguments: map[string]any{"lat": 0, "lon": 0},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError {
		t.Fatal("expected TOKEN_INVALID error")
	}
}
