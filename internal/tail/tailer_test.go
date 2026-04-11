package tail_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/tail"
)

func TestNewTailer_InvalidOptions(t *testing.T) {
	cases := []struct {
		name string
		opts tail.Options
	}{
		{"both zero", tail.Options{}},
		{"negative lines", tail.Options{MaxLines: -1}},
		{"negative bytes", tail.Options{MaxBytes: -1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tail.NewTailer(tc.opts)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNewTailer_ValidOptions(t *testing.T) {
	_, err := tail.NewTailer(tail.Options{MaxLines: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTailer_Apply_MaxLines(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	tr, _ := tail.NewTailer(tail.Options{MaxLines: 3})
	got, err := tr.Apply(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line3" || got[2] != "line5" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTailer_Apply_FewerLinesThanMax(t *testing.T) {
	input := "only\ntwo\n"
	tr, _ := tail.NewTailer(tail.Options{MaxLines: 10})
	got, err := tr.Apply(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestTailer_Apply_MaxBytes(t *testing.T) {
	// Each line is 5 chars + newline = 6 bytes. Request 13 bytes → 2 full lines.
	input := "aaaaa\nbbbbb\nccccc\n"
	tr, _ := tail.NewTailer(tail.Options{MaxBytes: 13})
	got, err := tr.Apply(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "bbbbb" || got[1] != "ccccc" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTailer_Apply_EmptyInput(t *testing.T) {
	tr, _ := tail.NewTailer(tail.Options{MaxLines: 5})
	got, err := tr.Apply(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
