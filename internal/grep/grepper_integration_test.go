package grep_test

import (
	"testing"

	"github.com/yourorg/fincut/internal/grep"
)

func TestGrepper_LargeInput_MatchCount(t *testing.T) {
	lines := make([]string, 1000)
	for i := range lines {
		if i%10 == 0 {
			lines[i] = "ERROR something went wrong"
		} else {
			lines[i] = "INFO all good"
		}
	}
	g, err := grep.New(grep.Options{Patterns: []string{"ERROR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	matches := g.Apply(lines)
	if len(matches) != 100 {
		t.Errorf("expected 100 matches, got %d", len(matches))
	}
}

func TestGrepper_ContextDoesNotExceedBounds(t *testing.T) {
	lines := []string{"ERROR first", "INFO second"}
	g, _ := grep.New(grep.Options{
		Patterns:      []string{"ERROR"},
		ContextBefore: 5,
		ContextAfter:  5,
	})
	matches := g.Apply(lines)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match")
	}
	if len(matches[0].Before) != 0 {
		t.Errorf("expected no before context, got %v", matches[0].Before)
	}
	if len(matches[0].After) != 1 {
		t.Errorf("expected 1 after line, got %v", matches[0].After)
	}
}

func TestGrepper_InvertWithContext(t *testing.T) {
	lines := []string{"INFO a", "ERROR b", "INFO c", "ERROR d", "INFO e"}
	g, _ := grep.New(grep.Options{
		Patterns:     []string{"ERROR"},
		Invert:       true,
		ContextAfter: 1,
	})
	matches := g.Apply(lines)
	// Inverted: INFO lines match (3 of them)
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
}
