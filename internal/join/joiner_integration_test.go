package join_test

import (
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/join"
)

func TestJoiner_LargeInput_CorrectGroupCount(t *testing.T) {
	const total = 1000
	const groupSize = 10

	lines := make([]string, total)
	for i := range lines {
		lines[i] = "line"
	}

	j, err := join.New(join.Options{GroupSize: groupSize, Separator: " "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := j.Apply(lines)
	expected := total / groupSize
	if len(out) != expected {
		t.Errorf("expected %d groups, got %d", expected, len(out))
	}
}

func TestJoiner_AllLinesPreserved(t *testing.T) {
	input := []string{"alpha", "beta", "gamma", "delta"}
	j, _ := join.New(join.Options{GroupSize: 0, Separator: "|"})
	out := j.Apply(input)

	if len(out) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(out))
	}

	parts := strings.Split(out[0], "|")
	if len(parts) != len(input) {
		t.Errorf("expected %d parts after split, got %d", len(input), len(parts))
	}
	for i, p := range parts {
		if p != input[i] {
			t.Errorf("part %d: expected %q, got %q", i, input[i], p)
		}
	}
}
