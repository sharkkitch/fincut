package timestamp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/user/fincut/internal/timestamp"
)

func TestStamper_LargeInput_AllLinesStamped(t *testing.T) {
	const n = 1000
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	s, err := timestamp.New(timestamp.Options{Prepend: true, Format: "2006"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(lines)
	if len(out) != n {
		t.Fatalf("expected %d lines, got %d", n, len(out))
	}
	for i, line := range out {
		if strings.TrimSpace(line) == "" {
			t.Errorf("line %d is blank", i)
		}
	}
}

func TestStamper_OriginalSliceUnmodified(t *testing.T) {
	original := []string{"alpha", "beta", "gamma"}
	copy := append([]string{}, original...)
	s, _ := timestamp.New(timestamp.Options{Append: true})
	s.Apply(original)
	for i, v := range original {
		if v != copy[i] {
			t.Errorf("original slice modified at index %d", i)
		}
	}
}
