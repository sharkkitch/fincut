package drop_test

import (
	"fmt"
	"testing"

	"github.com/your-org/fincut/internal/drop"
)

func TestDropper_LargeInput_CorrectCount(t *testing.T) {
	const total = 10_000
	lines := make([]string, total)
	for i := range lines {
		if i%3 == 0 {
			lines[i] = fmt.Sprintf("DEBUG line %d", i)
		} else {
			lines[i] = fmt.Sprintf("INFO line %d", i)
		}
	}

	d, err := drop.New(drop.Options{Patterns: []string{`^DEBUG`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := d.Apply(lines)

	// Every 3rd line (0, 3, 6, …) is DEBUG — ceil(10000/3) = 3334 lines dropped.
	expectedDropped := 0
	for i := 0; i < total; i++ {
		if i%3 == 0 {
			expectedDropped++
		}
	}
	expectedKept := total - expectedDropped

	if len(got) != expectedKept {
		t.Errorf("expected %d kept lines, got %d", expectedKept, len(got))
	}
}

func TestDropper_OriginalSliceUnmodified(t *testing.T) {
	input := []string{"DEBUG a", "INFO b", "DEBUG c"}
	copy := make([]string, len(input))
	copy[0], copy[1], copy[2] = input[0], input[1], input[2]

	d, _ := drop.New(drop.Options{Patterns: []string{`^DEBUG`}})
	_ = d.Apply(input)

	for i, line := range input {
		if line != copy[i] {
			t.Errorf("original slice modified at index %d: got %q, want %q", i, line, copy[i])
		}
	}
}
