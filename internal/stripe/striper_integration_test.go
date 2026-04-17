package stripe_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/fincut/internal/stripe"
)

func TestStriper_LargeInput_CorrectCount(t *testing.T) {
	const total = 1000
	lines := make([]string, total)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	s, err := stripe.New(stripe.Options{Every: 5, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(lines)
	want := total / 5
	if len(out) != want {
		t.Fatalf("expected %d lines, got %d", want, len(out))
	}
}

func TestStriper_OriginalSliceUnmodified(t *testing.T) {
	lines := []string{"x", "y", "z", "w"}
	copy_ := make([]string, len(lines))
	copy(copy_, lines)
	s, _ := stripe.New(stripe.Options{Every: 2, Offset: 0})
	s.Apply(lines)
	for i, v := range lines {
		if v != copy_[i] {
			t.Errorf("original slice modified at index %d", i)
		}
	}
}

func TestStriper_AllOffsets_CoverAllLines(t *testing.T) {
	const every = 3
	const total = 9
	lines := make([]string, total)
	for i := range lines {
		lines[i] = fmt.Sprintf("%d", i)
	}
	seen := make(map[string]bool)
	for off := 0; off < every; off++ {
		s, _ := stripe.New(stripe.Options{Every: every, Offset: off})
		for _, l := range s.Apply(lines) {
			seen[l] = true
		}
	}
	if len(seen) != total {
		t.Fatalf("expected all %d lines covered, got %d", total, len(seen))
	}
}
