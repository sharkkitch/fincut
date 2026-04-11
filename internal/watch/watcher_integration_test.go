package watch_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/user/fincut/internal/watch"
)

func TestWatcher_MultipleWrites(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "watch-multi-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w, err := watch.NewWatcher(watch.Options{
		Path:         f.Name(),
		Output:       &buf,
		PollInterval: 40 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	go w.Run(ctx) //nolint:errcheck

	lines := []string{"line one\n", "line two\n", "line three\n"}
	for _, l := range lines {
		time.Sleep(60 * time.Millisecond)
		f.WriteString(l)
	}

	<-ctx.Done()

	result := buf.String()
	for _, l := range lines {
		if !strings.Contains(result, strings.TrimSpace(l)) {
			t.Errorf("expected output to contain %q", l)
		}
	}
}

func TestWatcher_CancelStopsImmediately(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "watch-cancel-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w, _ := watch.NewWatcher(watch.Options{
		Path:         f.Name(),
		Output:       &buf,
		PollInterval: 50 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(300 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancel")
	}
}
