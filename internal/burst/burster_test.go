package burst

import (
	"testing"
	"time"
)

const (
	testPattern = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`
	testLayout  = "2006-01-02T15:04:05"
)

func defaultOpts() Options {
	return Options{
		TimestampPattern: testPattern,
		TimestampLayout:  testLayout,
		Window:           2 * time.Second,
		Threshold:        2.0,
	}
}

func TestNew_MissingPattern(t *testing.T) {
	opts := defaultOpts()
	opts.TimestampPattern = ""
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for empty TimestampPattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	opts := defaultOpts()
	opts.TimestampPattern = "["
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	opts := defaultOpts()
	opts.Window = 0
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestNew_ZeroThreshold(t *testing.T) {
	opts := defaultOpts()
	opts.Threshold = 0
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected error for zero Threshold")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	_, err := New(defaultOpts())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBurster_Apply_NoBurst(t *testing.T) {
	b, _ := New(defaultOpts())
	lines := []string{
		"2024-01-01T00:00:00 slow log line",
		"2024-01-01T00:00:10 another slow line",
	}
	bursts, err := b.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bursts) != 0 {
		t.Fatalf("expected 0 bursts, got %d", len(bursts))
	}
}

func TestBurster_Apply_DetectsBurst(t *testing.T) {
	b, _ := New(defaultOpts())
	lines := []string{
		"2024-01-01T00:00:00 line one",
		"2024-01-01T00:00:00 line two",
		"2024-01-01T00:00:01 line three",
		"2024-01-01T00:00:01 line four",
		"2024-01-01T00:00:02 line five",
	}
	bursts, err := b.Apply(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bursts) == 0 {
		t.Fatal("expected at least one burst")
	}
	if bursts[0].Rate <= 0 {
		t.Error("expected positive rate")
	}
}

func TestBurster_Apply_EmptyInput(t *testing.T) {
	b, _ := New(defaultOpts())
	bursts, err := b.Apply([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bursts) != 0 {
		t.Fatalf("expected 0 bursts for empty input, got %d", len(bursts))
	}
}
