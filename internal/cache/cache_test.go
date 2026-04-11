package cache_test

import (
	"testing"
	"time"

	"github.com/yourorg/fincut/internal/cache"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := cache.New(0, time.Second)
	if err != cache.ErrInvalidCapacity {
		t.Fatalf("expected ErrInvalidCapacity, got %v", err)
	}
}

func TestNew_InvalidTTL(t *testing.T) {
	_, err := cache.New(10, 0)
	if err != cache.ErrInvalidTTL {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c, err := cache.New(5, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	e := &cache.Entry{Lines: []string{"line1", "line2"}, ByteStart: 0, ByteEnd: 12}
	c.Set("file.log", e)

	got, ok := c.Get("file.log")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if len(got.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got.Lines))
	}
}

func TestCache_Get_Missing(t *testing.T) {
	c, _ := cache.New(5, time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for nonexistent key")
	}
}

func TestCache_Get_Expired(t *testing.T) {
	c, _ := cache.New(5, 10*time.Millisecond)
	c.Set("k", &cache.Entry{Lines: []string{"x"}})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestCache_Delete(t *testing.T) {
	c, _ := cache.New(5, time.Minute)
	c.Set("k", &cache.Entry{Lines: []string{"a"}})
	c.Delete("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestCache_EvictsOldestWhenFull(t *testing.T) {
	c, _ := cache.New(2, time.Minute)
	c.Set("first", &cache.Entry{Lines: []string{"a"}})
	time.Sleep(2 * time.Millisecond)
	c.Set("second", &cache.Entry{Lines: []string{"b"}})
	time.Sleep(2 * time.Millisecond)
	c.Set("third", &cache.Entry{Lines: []string{"c"}})

	if c.Len() != 2 {
		t.Fatalf("expected len 2 after eviction, got %d", c.Len())
	}
	_, ok := c.Get("first")
	if ok {
		t.Fatal("expected 'first' to have been evicted")
	}
}
