package dedup

import (
	"testing"
)

func TestNew_NegativeWindowSize(t *testing.T) {
	_, err := New(Options{WindowSize: -1})
	if err == nil {
		t.Fatal("expected error for negative window size")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	d, err := New(Options{WindowSize: 10, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduper")
	}
}

func TestDeduper_Apply_RemovesDuplicates(t *testing.T) {
	d, _ := New(Options{})
	input := []string{"alpha", "beta", "alpha", "gamma", "beta"}
	out := d.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique lines, got %d", len(out))
	}
}

func TestDeduper_Apply_AllUnique(t *testing.T) {
	d, _ := New(Options{})
	input := []string{"a", "b", "c"}
	out := d.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestDeduper_Apply_CaseInsensitive(t *testing.T) {
	d, _ := New(Options{CaseSensitive: false})
	input := []string{"Hello", "hello", "HELLO"}
	out := d.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 unique line (case-insensitive), got %d", len(out))
	}
}

func TestDeduper_Apply_CaseSensitive(t *testing.T) {
	d, _ := New(Options{CaseSensitive: true})
	input := []string{"Hello", "hello", "HELLO"}
	out := d.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique lines (case-sensitive), got %d", len(out))
	}
}

func TestDeduper_Apply_WindowEviction(t *testing.T) {
	d, _ := New(Options{WindowSize: 2})
	// After window evicts "alpha", it should be seen as new again
	out := d.Apply([]string{"alpha", "beta", "gamma", "alpha"})
	if len(out) != 4 {
		t.Fatalf("expected 4 lines after window eviction, got %d: %v", len(out), out)
	}
}

func TestDeduper_Reset_ClearsState(t *testing.T) {
	d, _ := New(Options{})
	d.Apply([]string{"foo", "bar"})
	d.Reset()
	out := d.Apply([]string{"foo", "bar"})
	if len(out) != 2 {
		t.Fatalf("expected 2 lines after reset, got %d", len(out))
	}
}

func TestDeduper_Apply_EmptyInput(t *testing.T) {
	d, _ := New(Options{})
	out := d.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}
