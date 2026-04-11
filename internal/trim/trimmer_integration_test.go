package trim_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/trim"
)

// TestTrimmer_CombinedOptions verifies that MaxLines and StripEmpty
// interact correctly: empty lines are stripped before the line cap is applied.
func TestTrimmer_CombinedOptions(t *testing.T) {
	tr, err := trim.NewTrimmer(trim.Options{
		MaxLines:   3,
		StripEmpty: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"a", "", "b", "", "c", "d"}
	var buf strings.Builder
	tr.Apply(lines, &buf)
	result := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(result) != 3 {
		t.Errorf("expected 3 lines after strip+cap, got %d: %v", len(result), result)
	}
	if result[0] != "a" || result[1] != "b" || result[2] != "c" {
		t.Errorf("unexpected content: %v", result)
	}
}

// TestTrimmer_BytesAndLines ensures the more restrictive limit wins.
func TestTrimmer_BytesAndLines(t *testing.T) {
	// 4 lines of "hello" = 4*(5+1)=24 bytes; MaxBytes=12 allows only 2
	tr, err := trim.NewTrimmer(trim.Options{
		MaxLines: 10,
		MaxBytes: 12,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"hello", "hello", "hello", "hello"}
	var buf strings.Builder
	tr.Apply(lines, &buf)
	result := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(result) != 2 {
		t.Errorf("expected 2 lines (byte limit), got %d", len(result))
	}
}

// TestTrimmer_EmptyInput verifies graceful handling of no lines.
func TestTrimmer_EmptyInput(t *testing.T) {
	tr, _ := trim.NewTrimmer(trim.Options{MaxLines: 5})
	var buf strings.Builder
	n, err := tr.Apply([]string{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 bytes written, got %d", n)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}
