package window

import (
	"testing"
)

func lines(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = string(rune('A' + i))
	}
	return out
}

func TestNew_ZeroSize(t *testing.T) {
	_, err := New(Options{Size: 0})
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestNew_NegativeStep(t *testing.T) {
	_, err := New(Options{Size: 3, Step: -1})
	if err == nil {
		t.Fatal("expected error for negative step")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	w, err := New(Options{Size: 4, Step: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.size != 4 || w.step != 2 {
		t.Fatalf("unexpected fields: size=%d step=%d", w.size, w.step)
	}
}

func TestNew_DefaultStep(t *testing.T) {
	w, err := New(Options{Size: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.step != 5 {
		t.Fatalf("expected step=5, got %d", w.step)
	}
}

func TestWindower_Apply_NonOverlapping(t *testing.T) {
	w, _ := New(Options{Size: 3})
	win := w.Apply(lines(9))
	if len(win) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(win))
	}
	for _, ww := range win {
		if len(ww) != 3 {
			t.Fatalf("expected window size 3, got %d", len(ww))
		}
	}
}

func TestWindower_Apply_Overlapping(t *testing.T) {
	w, _ := New(Options{Size: 4, Step: 2})
	win := w.Apply(lines(8))
	// starts: 0,2,4,6 -> 4 windows
	if len(win) != 4 {
		t.Fatalf("expected 4 windows, got %d", len(win))
	}
}

func TestWindower_Apply_EmptyInput(t *testing.T) {
	w, _ := New(Options{Size: 3})
	win := w.Apply(nil)
	if win != nil {
		t.Fatalf("expected nil for empty input, got %v", win)
	}
}

func TestWindower_Apply_PartialLastWindow(t *testing.T) {
	w, _ := New(Options{Size: 4})
	win := w.Apply(lines(6))
	// windows: [0-3], [4-5]
	if len(win) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(win))
	}
	if len(win[1]) != 2 {
		t.Fatalf("expected last window size 2, got %d", len(win[1]))
	}
}

func TestFlatten_NonOverlapping(t *testing.T) {
	w, _ := New(Options{Size: 3})
	src := lines(9)
	win := w.Apply(src)
	flat := Flatten(win)
	if len(flat) != 9 {
		t.Fatalf("expected 9 lines, got %d", len(flat))
	}
}
