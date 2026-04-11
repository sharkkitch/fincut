package watch

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"
)

func TestNewWatcher_EmptyPath(t *testing.T) {
	_, err := NewWatcher(Options{Output: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewWatcher_NilOutput(t *testing.T) {
	_, err := NewWatcher(Options{Path: "file.log"})
	if err == nil {
		t.Fatal("expected error for nil output")
	}
}

func TestNewWatcher_DefaultInterval(t *testing.T) {
	w, err := NewWatcher(Options{Path: "file.log", Output: &bytes.Buffer{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.interval != 500*time.Millisecond {
		t.Errorf("expected default interval 500ms, got %v", w.interval)
	}
}

func TestNewWatcher_CustomInterval(t *testing.T) {
	w, err := NewWatcher(Options{
		Path:         "file.log",
		Output:       &bytes.Buffer{},
		PollInterval: 200 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.interval != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", w.interval)
	}
}

func TestWatcher_Run_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	w, _ := NewWatcher(Options{Path: "/nonexistent/file.log", Output: &buf})
	ctx := context.Background()
	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatcher_Run_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "watch-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w, err := NewWatcher(Options{
		Path:         f.Name(),
		Output:       &buf,
		PollInterval: 50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	time.Sleep(80 * time.Millisecond)
	f.WriteString("hello watch\n")

	<-done

	if !bytes.Contains(buf.Bytes(), []byte("hello watch")) {
		t.Errorf("expected output to contain 'hello watch', got: %q", buf.String())
	}
}
