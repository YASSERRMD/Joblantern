package main

import (
	"context"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/yasserrmd/joblantern/internal/registry"
)

type fakeProv struct {
	matches []registry.Match
	get     *registry.Company
	err     error
}

func (f *fakeProv) Name() string { return "fake" }
func (f *fakeProv) LookupByName(_ context.Context, _, _ string, _ int) ([]registry.Match, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.matches, nil
}
func (f *fakeProv) Get(_ context.Context, _ string) (*registry.Company, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.get, nil
}

func TestLookupCompany_OK(t *testing.T) {
	fp := &fakeProv{matches: []registry.Match{{ID: "gb/1", Name: "Acme"}}}
	srv := newServer(fp)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, _ := srv.Connect(ctx, st, nil)
	defer func() { _ = ss.Close() }()
	c := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, _ := c.Connect(ctx, ct, nil)
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "lookup_company", Arguments: map[string]any{"name": "Acme"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}

func TestCheckRegistrationStatus_Recent(t *testing.T) {
	fp := &fakeProv{get: &registry.Company{
		Match: registry.Match{
			ID: "gb/1", Name: "Fresh",
			Status:            "Active",
			IncorporationDate: time.Now().Add(-60 * 24 * time.Hour),
		},
	}}
	srv := newServer(fp)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, _ := srv.Connect(ctx, st, nil)
	defer func() { _ = ss.Close() }()
	c := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, _ := c.Connect(ctx, ct, nil)
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "check_registration_status", Arguments: map[string]any{"id": "gb/1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}
