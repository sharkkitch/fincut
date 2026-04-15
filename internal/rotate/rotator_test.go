package rotate

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestNewRotator_EmptyPath(t *testing.T) {
	_, err := NewRotator(Options{Output: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewRotator_NilOutput(t *testing.T) {
	_, err := NewRotator(Options{Path: "/tmp/test.log"})
	if err == nil {
		t.Fatal("expected error for nil output")
	}
}

func TestNewRotator_DefaultInterval(t *testing.T) {
	r, err := NewRotator(Options{Path: "/tmp/test.log", Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.opts.Interval != 5*time.Second {
		t.Errorf("expected default interval 5s, got %v", r.opts.Interval)
	}
}

func TestNewRotator_CustomInterval(t *testing.T) {
	r, err := NewRotator(Options{
		Path:     "/tmp/test.log",
		Output:   &bytes.Buffer{},
		Interval: 2 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.opts.Interval != 2*time.Second {
		t.Errorf("expected 2s interval, got %v", r.opts.Interval)
	}
}

func TestRotator_Detect_NoRotation(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rottest")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello\n")
	f.Close()

	r, err := NewRotator(Options{Path: f.Name(), Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := r.Baseline()
	if err != nil {
		t.Fatal(err)
	}
	rotated, err := r.Detect(baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rotated {
		t.Error("expected no rotation detected")
	}
}

func TestRotator_Detect_Truncation(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "rottest")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello world\n")
	f.Close()

	r, err := NewRotator(Options{Path: f.Name(), Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := r.Baseline()
	if err != nil {
		t.Fatal(err)
	}

	// Truncate the file to simulate rotation.
	if err := os.WriteFile(f.Name(), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	rotated, err := r.Detect(baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rotated {
		t.Error("expected rotation detected after truncation")
	}
}

func TestRotator_Detect_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/missing.log"

	r, err := NewRotator(Options{Path: path, Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatal(err)
	}
	rotated, err := r.Detect(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nil baseline — no rotation yet
	if rotated {
		t.Error("expected no rotation with nil baseline")
	}
}
