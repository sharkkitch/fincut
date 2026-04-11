package rotate

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestNewRotator_EmptyPath(t *testing.T) {
	_, err := NewRotator(Options{Path: "", Output: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewRotator_NilOutput(t *testing.T) {
	_, err := NewRotator(Options{Path: "/tmp/test.log", Output: nil})
	if err == nil {
		t.Fatal("expected error for nil output")
	}
}

func TestNewRotator_DefaultInterval(t *testing.T) {
	r, err := NewRotator(Options{Path: "/tmp/test.log", Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Interval() != 2*time.Second {
		t.Errorf("expected default interval 2s, got %v", r.Interval())
	}
}

func TestNewRotator_CustomInterval(t *testing.T) {
	r, err := NewRotator(Options{Path: "/tmp/test.log", Output: &bytes.Buffer{}, Interval: 5 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Interval() != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", r.Interval())
	}
}

func TestRotator_Detect_NoRotation(t *testing.T) {
	f, err := os.CreateTemp("", "rotator-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("initial content\n")
	f.Close()

	r, err := NewRotator(Options{Path: f.Name(), Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Snapshot(); err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}
	rotated, err := r.Detect()
	if err != nil {
		t.Fatalf("detect failed: %v", err)
	}
	if rotated {
		t.Error("expected no rotation detected")
	}
}

func TestRotator_Detect_Truncation(t *testing.T) {
	f, err := os.CreateTemp("", "rotator-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("some log content that is long enough\n")
	f.Close()

	r, err := NewRotator(Options{Path: f.Name(), Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Snapshot(); err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}

	// Truncate the file to simulate rotation
	if err := os.Truncate(f.Name(), 0); err != nil {
		t.Fatalf("truncate failed: %v", err)
	}

	rotated, err := r.Detect()
	if err != nil {
		t.Fatalf("detect failed: %v", err)
	}
	if !rotated {
		t.Error("expected rotation to be detected after truncation")
	}
}
