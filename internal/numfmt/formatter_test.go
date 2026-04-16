package numfmt

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(Options{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_NoCaptureGroup(t *testing.T) {
	_, err := New(Options{Pattern: `\d+`})
	if err == nil {
		t.Fatal("expected error when no capture group")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	f, err := New(Options{Pattern: `=(\d+)`, Precision: 2, Thousands: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestFormatter_Apply_Integer(t *testing.T) {
	f, _ := New(Options{Pattern: `val=(\d+)`, Thousands: false})
	out := f.Apply([]string{"val=42 end", "no match"})
	if out[0] != "val=42 end" {
		t.Errorf("unexpected: %q", out[0])
	}
	if out[1] != "no match" {
		t.Errorf("unexpected: %q", out[1])
	}
}

func TestFormatter_Apply_Thousands(t *testing.T) {
	f, _ := New(Options{Pattern: `n=(\d+)`, Thousands: true})
	out := f.Apply([]string{"n=1000000"})
	if out[0] != "n=1,000,000" {
		t.Errorf("expected n=1,000,000, got %q", out[0])
	}
}

func TestFormatter_Apply_Float(t *testing.T) {
	f, _ := New(Options{Pattern: `v=([0-9]+\.[0-9]+)`, Precision: 2})
	out := f.Apply([]string{"v=3.14159"})
	if out[0] != "v=3.14" {
		t.Errorf("expected v=3.14, got %q", out[0])
	}
}

func TestFormatter_Apply_EmptyInput(t *testing.T) {
	f, _ := New(Options{Pattern: `(\d+)`})
	out := f.Apply([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output")
	}
}

func TestFormatter_Apply_NoMatch(t *testing.T) {
	f, _ := New(Options{Pattern: `val=(\d+)`, Thousands: true})
	out := f.Apply([]string{"nothing here"})
	if out[0] != "nothing here" {
		t.Errorf("line should be unchanged: %q", out[0])
	}
}

func TestFormatInt_Small(t *testing.T) {
	result := formatInt(999, true)
	if result != "999" {
		t.Errorf("expected 999, got %q", result)
	}
}

func TestFormatInt_Negative(t *testing.T) {
	result := formatInt(-1234567, true)
	if result != "-1,234,567" {
		t.Errorf("expected -1,234,567, got %q", result)
	}
}
