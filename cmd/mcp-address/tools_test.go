package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/cache"
	"github.com/yasserrmd/joblantern/internal/nominatim"
	"github.com/yasserrmd/joblantern/internal/overpass"
)

func newTestServer(nomURL, overURL string) *mcp.Server {
	d := deps{
		nom:   nominatim.New(nomURL),
		over:  overpass.New(overURL),
		cache: cache.New[string, any](time.Minute),
	}
	return newServer(d)
}

func TestVerifyAddressExists_OK(t *testing.T) {
	nom := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"place_id":1,"osm_type":"node","osm_id":1,"lat":"25.2","lon":"55.27","display_name":"Burj Khalifa, Dubai"}]`))
	}))
	defer nom.Close()
	srv := newTestServer(nom.URL, "")

	cli, ssn := connect(t, srv)
	defer cli.Close()
	defer ssn.Close()

	res, err := cli.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "verify_address_exists", Arguments: map[string]any{"address": "Burj Khalifa, Dubai"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}

func TestClassifyBuildingType_Commercial(t *testing.T) {
	over := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"elements":[{"type":"way","id":1,"tags":{"office":"company","building":"office"}}]}`))
	}))
	defer over.Close()
	srv := newTestServer("", over.URL)

	cli, ssn := connect(t, srv)
	defer cli.Close()
	defer ssn.Close()

	res, err := cli.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "classify_building_type",
		Arguments: map[string]any{"lat": 25.2, "lon": 55.27, "radius_m": 50},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}

func connect(t *testing.T, srv *mcp.Server) (*mcp.ClientSession, *mcp.ServerSession) {
	t.Helper()
	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, err := srv.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	c := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, err := c.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	return cs, ss
}
