package cache_test

import (
	"testing"
	"time"

	"github.com/yasserrmd/joblantern/internal/cache"
)

func TestTTL_GetSet(t *testing.T) {
	c := cache.New[string, int](100 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected miss")
	}
	c.Set("a", 42)
	if v, ok := c.Get("a"); !ok || v != 42 {
		t.Fatalf("got %v ok=%v want 42 true", v, ok)
	}
	time.Sleep(150 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected expiry")
	}
}
