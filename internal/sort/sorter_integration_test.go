package sort_test

import (
	"fmt"
	"testing"

	fincsort "github.com/user/fincut/internal/sort"
)

func TestSorter_LargeInput_PreservesCount(t *testing.T) {
	const n = 10_000
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = fmt.Sprintf("line-%05d", n-i)
	}
	s, err := fincsort.New(fincsort.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(lines)
	if len(out) != n {
		t.Errorf("expected %d lines, got %d", n, len(out))
	}
	if out[0] != "line-00001" {
		t.Errorf("expected line-00001 first, got %q", out[0])
	}
}

func TestSorter_UniqueOnLargeInput(t *testing.T) {
	const n = 1_000
	lines := make([]string, n*2)
	for i := 0; i < n; i++ {
		lines[i] = fmt.Sprintf("key-%04d", i)
		lines[i+n] = fmt.Sprintf("key-%04d", i)
	}
	s, err := fincsort.New(fincsort.Options{Unique: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(lines)
	if len(out) != n {
		t.Errorf("expected %d unique lines, got %d", n, len(out))
	}
}

func TestSorter_OriginalSliceUnmodified(t *testing.T) {
	in := []string{"z", "a", "m"}
	orig := []string{"z", "a", "m"}
	s, _ := fincsort.New(fincsort.Options{})
	s.Apply(in)
	for i, v := range orig {
		if in[i] != v {
			t.Errorf("original slice modified at %d: got %q, want %q", i, in[i], v)
		}
	}
}
