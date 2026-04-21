package rotate_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/yourorg/fincut/internal/rotate"
)

func writeTempRotate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rotate-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNewRotator_EmptyPath(t *testing.T) {
	_, err := rotate.NewRotator(rotate.Options{Out: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewRotator_NilOutput(t *testing.T) {
	path := writeTempRotate(t, "hello\n")
	_, err := rotate.NewRotator(rotate.Options{Path: path})
	if err == nil {
		t.Fatal("expected error for nil output")
	}
}

func TestNewRotator_DefaultInterval(t *testing.T) {
	path := writeTempRotate(t, "line\n")
	buf := &bytes.Buffer{}
	r, err := rotate.NewRotator(rotate.Options{Path: path, Out: buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rotator")
	}
}

func TestNewRotator_CustomInterval(t *testing.T) {
	path := writeTempRotate(t, "line\n")
	buf := &bytes.Buffer{}
	_, err := rotate.NewRotator(rotate.Options{
		Path:     path,
		Out:      buf,
		Interval: 500 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRotator_Detect_NoRotation(t *testing.T) {
	path := writeTempRotate(t, "stable content\n")
	buf := &bytes.Buffer{}
	r, err := rotate.NewRotator(rotate.Options{Path: path, Out: buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rotated, err := r.Detect()
	if err != nil {
		t.Fatalf("detect error: %v", err)
	}
	if rotated {
		t.Error("expected no rotation for stable file")
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestRotator_Detect_AfterTruncate(t *testing.T) {
	path := writeTempRotate(t, "initial content with many bytes\n")
	buf := &bytes.Buffer{}
	r, err := rotate.NewRotator(rotate.Options{Path: path, Out: buf})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	// Truncate the file to simulate rotation
	if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	rotated, err := r.Detect()
	if err != nil {
		t.Fatalf("detect error: %v", err)
	}
	if !rotated {
		t.Error("expected rotation to be detected after truncation")
	}
	if buf.Len() == 0 {
		t.Error("expected event written to output")
	}
}
