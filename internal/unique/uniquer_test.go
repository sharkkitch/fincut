package unique

import (
	"testing"
)

func TestNew_NegativeField(t *testing.T) {
	_, err := New(Options{Field: -1, Delimiter: ","})
	if err == nil {
		t.Fatal("expected error for negative field")
	}
}

func TestNew_FieldWithoutDelimiter(t *testing.T) {
	_, err := New(Options{Field: 2})
	if err == nil {
		t.Fatal("expected error for field without delimiter")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = New(Options{Field: 1, Delimiter: ":"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUniquer_Apply_WholeLineKey(t *testing.T) {
	u, _ := New(Options{})
	input := []string{"alpha", "beta", "alpha", "gamma", "beta"}
	got := u.Apply(input)
	want := []string{"alpha", "beta", "gamma"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, got[i])
		}
	}
}

func TestUniquer_Apply_ByField(t *testing.T) {
	u, _ := New(Options{Field: 2, Delimiter: ","})
	input := []string{"a,foo", "b,bar", "c,foo", "d,baz"}
	got := u.Apply(input)
	want := []string{"a,foo", "b,bar", "d,baz"}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: want %q, got %q", i, w, got[i])
		}
	}
}

func TestUniquer_Apply_CaseInsensitive(t *testing.T) {
	u, _ := New(Options{CaseInsensitive: true})
	input := []string{"Hello", "HELLO", "world", "World"}
	got := u.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 unique lines, got %d", len(got))
	}
}

func TestUniquer_Apply_EmptyInput(t *testing.T) {
	u, _ := New(Options{})
	got := u.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %d lines", len(got))
	}
}

func TestUniquer_Apply_FieldOutOfRange(t *testing.T) {
	u, _ := New(Options{Field: 5, Delimiter: ","})
	input := []string{"a,b", "a,b", "c,d"}
	// field 5 is out of range — falls back to whole line as key
	got := u.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}
