package skip

import (
	"testing"
)

func TestNew_EveryTooSmall(t *testing.T) {
	_, err := New(Options{Every: 1})
	if err == nil {
		t.Fatal("expected error for Every=1")
	}
}

func TestNew_NegativeOffset(t *testing.T) {
	_, err := New(Options{Every: 2, Offset: -1})
	if err == nil {
		t.Fatal("expected error for negative Offset")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	s, err := New(Options{Every: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Skipper")
	}
}

func TestSkipper_Apply_EmptyInput(t *testing.T) {
	s, _ := New(Options{Every: 2})
	out := s.Apply(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(out))
	}
}

func TestSkipper_Apply_EverySecond(t *testing.T) {
	s, _ := New(Options{Every: 2})
	in := []string{"a", "b", "c", "d", "e", "f"}
	out := s.Apply(in)
	// lines at positions 2,4,6 (b,d,f) should be dropped
	want := []string{"a", "c", "e"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestSkipper_Apply_EveryThird(t *testing.T) {
	s, _ := New(Options{Every: 3})
	in := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	out := s.Apply(in)
	// positions 3,6,9 (c,f,i) dropped
	want := []string{"a", "b", "d", "e", "g", "h"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestSkipper_Apply_WithOffset(t *testing.T) {
	// Offset=1 shifts window: effective position for index 0 is 0 (not dropped),
	// index 1 -> pos 1 (not dropped), index 2 -> pos 2 (dropped for Every=2), etc.
	s, _ := New(Options{Every: 2, Offset: 1})
	in := []string{"a", "b", "c", "d", "e"}
	out := s.Apply(in)
	// pos: a->0(skip<1), b->1, c->2(drop), d->3, e->4(drop)
	want := []string{"a", "b", "d"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(out), out)
	}
	for i, w := range want {
		if out[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, out[i])
		}
	}
}

func TestFormatSummary(t *testing.T) {
	got := FormatSummary(10, 7)
	if got == "" {
		t.Fatal("expected non-empty summary")
	}
}
