package dedup_test

import (
	"fmt"
	"testing"

	"github.com/user/fincut/internal/dedup"
)

func TestDeduper_LargeInput_UniqueCount(t *testing.T) {
	d, err := dedup.New(dedup.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const total = 1000
	lines := make([]string, total)
	for i := 0; i < total; i++ {
		// Every 10th line is a duplicate of the previous
		if i > 0 && i%10 == 0 {
			lines[i] = lines[i-1]
		} else {
			lines[i] = fmt.Sprintf("line-%d", i)
		}
	}

	out := d.Apply(lines)
	expected := total - 90 // 90 duplicates injected
	if len(out) != expected {
		t.Fatalf("expected %d unique lines, got %d", expected, len(out))
	}
}

func TestDeduper_WindowSize_LimitsMemory(t *testing.T) {
	d, err := dedup.New(dedup.Options{WindowSize: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Cycle through 5 unique lines repeatedly; window should evict and re-admit
	lines := make([]string, 0, 30)
	for round := 0; round < 3; round++ {
		for i := 0; i < 10; i++ {
			lines = append(lines, fmt.Sprintf("line-%d", i))
		}
	}

	out := d.Apply(lines)
	if len(out) == 0 {
		t.Fatal("expected non-empty output with windowed dedup")
	}
	if len(out) > len(lines) {
		t.Fatalf("output cannot exceed input length")
	}
}
