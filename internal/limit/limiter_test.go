package limit

import (
	"strings"
	"testing"
)

func TestNew_NeitherSet(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Fatal("expected error when neither MaxLines nor MaxBytes is set")
	}
}

func TestNew_NegativeMaxLines(t *testing.T) {
	_, err := New(Options{MaxLines: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxLines")
	}
}

func TestNew_NegativeMaxBytes(t *testing.T) {
	_, err := New(Options{MaxBytes: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxBytes")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	for _, opts := range []Options{
		{MaxLines: 10},
		{MaxBytes: 512},
		{MaxLines: 5, MaxBytes: 256},
	} {
		if _, err := New(opts); err != nil {
			t.Errorf("unexpected error for %+v: %v", opts, err)
		}
	}
}

func TestLimiter_Apply_MaxLines(t *testing.T) {
	lim, _ := New(Options{MaxLines: 3})
	input := []string{"a", "b", "c", "d", "e"}
	got := lim.Apply(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestLimiter_Apply_MaxBytes(t *testing.T) {
	// each line is 3 bytes ("abc"), limit to 9 bytes → 3 lines
	lim, _ := New(Options{MaxBytes: 9})
	input := []string{"abc", "abc", "abc", "abc"}
	got := lim.Apply(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestLimiter_Apply_BothLimits_LinesWins(t *testing.T) {
	// MaxLines=2 kicks in before MaxBytes=1000
	lim, _ := New(Options{MaxLines: 2, MaxBytes: 1000})
	input := []string{"hello", "world", "extra"}
	got := lim.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestLimiter_Apply_EmptyInput(t *testing.T) {
	lim, _ := New(Options{MaxLines: 10})
	got := lim.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(got))
	}
}

func TestLimiter_FormatSummary(t *testing.T) {
	cases := []struct {
		opts Options
		want string
	}{
		{Options{MaxLines: 5}, "limit: 5 lines"},
		{Options{MaxBytes: 128}, "limit: 128 bytes"},
		{Options{MaxLines: 3, MaxBytes: 64}, "limit: 3 lines / 64 bytes"},
	}
	for _, tc := range cases {
		lim, _ := New(tc.opts)
		got := lim.FormatSummary()
		if !strings.Contains(got, strings.Split(tc.want, " ")[1]) {
			t.Errorf("FormatSummary() = %q, want to contain %q", got, tc.want)
		}
	}
}
