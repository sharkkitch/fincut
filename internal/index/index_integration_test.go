package index_test

import (
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/index"
)

func TestIndex_LargeInput_OffsetConsistency(t *testing.T) {
	var sb strings.Builder
	lines := 1000
	for i := 0; i < lines; i++ {
		sb.WriteString("log line content here\n")
	}

	idx, err := index.Build(strings.NewReader(sb.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len() != lines {
		t.Fatalf("expected %d lines, got %d", lines, idx.Len())
	}

	// Each line is "log line content here\n" = 22 bytes
	for i := 0; i < lines; i++ {
		e, err := idx.Lookup(i)
		if err != nil {
			t.Fatalf("lookup failed at line %d: %v", i, err)
		}
		expected := int64(i * 22)
		if e.Offset != expected {
			t.Errorf("line %d: expected offset %d, got %d", i, expected, e.Offset)
		}
	}
}

func TestIndex_Range_CoversAllLines(t *testing.T) {
	input := "one\ntwo\nthree\nfour\nfive\n"
	idx, err := index.Build(strings.NewReader(input))
	if err != nil {
		t.Fatalf("build error: %v", err)
	}

	entries, err := idx.Range(0, idx.Len()-1)
	if err != nil {
		t.Fatalf("range error: %v", err)
	}
	if len(entries) != idx.Len() {
		t.Errorf("expected %d entries from full range, got %d", idx.Len(), len(entries))
	}
}
