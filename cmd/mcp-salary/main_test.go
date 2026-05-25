package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSalaryRangeCheck_ImplausiblyHigh(t *testing.T) {
	p, err := loadPack()
	if err != nil {
		t.Fatal(err)
	}
	srv := newServer(p)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()
	ss, _ := srv.Connect(ctx, st, nil)
	defer func() { _ = ss.Close() }()
	c := mcp.NewClient(&mcp.Implementation{Name: "t"}, nil)
	cs, _ := c.Connect(ctx, ct, nil)
	defer func() { _ = cs.Close() }()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name: "salary_range_check",
		Arguments: map[string]any{
			"country": "BD", "role": "warehouse_worker",
			"currency": "USD", "claimed_amount": 12000.0,
			"claimed_period": "month",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error: %v", res.Content)
	}
}
