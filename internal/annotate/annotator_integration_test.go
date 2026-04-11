package annotate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/user/fincut/internal/annotate"
)

func TestAnnotator_LargeInput_LineNumbersAreSequential(t *testing.T) {
	const n = 1000
	input := make([]string, n)
	for i := range input {
		input[i] = fmt.Sprintf("line content %d", i)
	}

	a, err := annotate.New(annotate.Options{LineNumbers: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := a.Apply(input)
	if len(out) != n {
		t.Fatalf("expected %d lines, got %d", n, len(out))
	}

	for i, line := range out {
		expectedNum := fmt.Sprintf("%d", i+1)
		if !strings.HasPrefix(line, expectedNum) {
			t.Errorf("line %d: expected prefix %q, got %q", i+1, expectedNum, line)
		}
	}
}

func TestAnnotator_PrefixAndLineNumbers_AllLinesAnnotated(t *testing.T) {
	input := []string{"alpha", "beta", "gamma", "delta"}
	a, err := annotate.New(annotate.Options{
		LineNumbers: true,
		Prefix:      "audit",
		Separator:   " | ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := a.Apply(input)
	for i, line := range out {
		if !strings.Contains(line, "audit") {
			t.Errorf("line %d missing prefix 'audit': %q", i+1, line)
		}
		if !strings.Contains(line, input[i]) {
			t.Errorf("line %d missing original content %q: %q", i+1, input[i], line)
		}
	}
}
