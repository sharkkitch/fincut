package limit_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/fincut/internal/limit"
)

func generateLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line-%05d", i)
	}
	return lines
}

func TestLimiter_LargeInput_MaxLines(t *testing.T) {
	const total = 100_000
	const cap = 250
	lim, err := limit.New(limit.Options{MaxLines: cap})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := generateLines(total)
	got := lim.Apply(input)
	if len(got) != cap {
		t.Fatalf("expected %d lines, got %d", cap, len(got))
	}
	if got[0] != input[0] || got[cap-1] != input[cap-1] {
		t.Error("output lines do not match the leading input lines")
	}
}

func TestLimiter_LargeInput_MaxBytes(t *testing.T) {
	// each line is "line-NNNNN" = 10 bytes
	const lineLen = 10
	const byteLimit = 1000
	expected := byteLimit / lineLen

	lim, err := limit.New(limit.Options{MaxBytes: byteLimit})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	input := generateLines(10_000)
	got := lim.Apply(input)
	if len(got) != expected {
		t.Fatalf("expected %d lines, got %d", expected, len(got))
	}
}

func TestLimiter_OriginalSliceUnmodified(t *testing.T) {
	lim, _ := limit.New(limit.Options{MaxLines: 5})
	input := generateLines(20)
	copy := append([]string(nil), input...)
	lim.Apply(input)
	for i, v := range copy {
		if input[i] != v {
			t.Errorf("input[%d] was modified: got %q want %q", i, input[i], v)
		}
	}
}
