package chunk

import (
	"strings"
	"testing"
)

func TestNew_MutuallyExclusive(t *testing.T) {
	_, err := New(Options{Size: 3, Delimiter: `---`})
	if err == nil {
		t.Fatal("expected error for mutually exclusive options")
	}
}

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither Size nor Delimiter is set")
	}
}

func TestNew_NegativeSize(t *testing.T) {
	_, err := New(Options{Size: -1})
	if err == nil {
		t.Fatal("expected error for negative Size")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Delimiter: `[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid delimiter pattern")
	}
}

func TestNew_DefaultLabelPrefix(t *testing.T) {
	c, err := New(Options{Size: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.opts.LabelPrefix != "chunk" {
		t.Errorf("expected default label prefix 'chunk', got %q", c.opts.LabelPrefix)
	}
}

func TestChunker_Apply_EmptyInput(t *testing.T) {
	c, _ := New(Options{Size: 3})
	if got := c.Apply(nil); got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}

func TestChunker_Apply_BySize(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	c, _ := New(Options{Size: 2})
	chunks := c.Apply(lines)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if len(chunks[2].Lines) != 1 {
		t.Errorf("last chunk should have 1 line, got %d", len(chunks[2].Lines))
	}
	if chunks[0].Label != "chunk-1" {
		t.Errorf("unexpected label: %s", chunks[0].Label)
	}
}

func TestChunker_Apply_BySizeExact(t *testing.T) {
	lines := []string{"a", "b", "c", "d"}
	c, _ := New(Options{Size: 2, LabelPrefix: "block"})
	chunks := c.Apply(lines)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if chunks[1].Label != "block-2" {
		t.Errorf("unexpected label: %s", chunks[1].Label)
	}
}

func TestChunker_Apply_ByDelimiter(t *testing.T) {
	lines := []string{"START", "a", "b", "START", "c"}
	c, _ := New(Options{Delimiter: `^START`})
	chunks := c.Apply(lines)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if !strings.HasPrefix(chunks[0].Lines[0], "START") {
		t.Errorf("first line of chunk should be delimiter line")
	}
}

func TestChunker_Apply_ByDelimiter_NoMatch(t *testing.T) {
	lines := []string{"a", "b", "c"}
	c, _ := New(Options{Delimiter: `^---`})
	chunks := c.Apply(lines)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk when delimiter never matches, got %d", len(chunks))
	}
	if len(chunks[0].Lines) != 3 {
		t.Errorf("expected all 3 lines in single chunk")
	}
}
