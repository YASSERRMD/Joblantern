package db

import (
	"context"
	"testing"
)

func TestTenantIsolation(t *testing.T) {
	ctx := WithTenant(context.Background(), "tenant-a")
	got, err := FromContext(ctx)
	if err != nil || got != "tenant-a" {
		t.Fatalf("expected tenant-a, got %q err=%v", got, err)
	}
	if _, err := FromContext(context.Background()); err == nil {
		t.Fatalf("expected error when no tenant set")
	}
}

func TestThreeTenantsThenOffboard(t *testing.T) {
	tenants := []string{"a", "b", "c"}
	for _, tn := range tenants {
		ctx := WithTenant(context.Background(), tn)
		got, err := FromContext(ctx)
		if err != nil || got != tn {
			t.Fatalf("tenant %s: got %q err=%v", tn, got, err)
		}
	}
}
