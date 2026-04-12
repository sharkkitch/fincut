package normalize_test

import (
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/normalize"
)

func TestNormalizer_LargeInput_AllLinesProcessed(t *testing.T) {
	n, err := normalize.New(normalize.Options{
		TrimSpace:      true,
		CollapseSpaces: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	const count = 10_000
	lines := make([]string, count)
	for i := range lines {
		lines[i] = "  log   entry   line  "
	}

	result := n.Apply(lines)
	if len(result) != count {
		t.Fatalf("expected %d lines, got %d", count, len(result))
	}
	for i, line := range result {
		if strings.HasPrefix(line, " ") || strings.HasSuffix(line, " ") {
			t.Errorf("line %d still has surrounding spaces: %q", i, line)
		}
		if strings.Contains(line, "  ") {
			t.Errorf("line %d still has consecutive spaces: %q", i, line)
		}
	}
}

func TestNormalizer_StripControl_PreservesTab(t *testing.T) {
	n, err := normalize.New(normalize.Options{StripControl: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := []string{"col1\tcol2\tcol3"}
	result := n.Apply(input)
	if result[0] != "col1\tcol2\tcol3" {
		t.Errorf("tab should be preserved, got %q", result[0])
	}
}
