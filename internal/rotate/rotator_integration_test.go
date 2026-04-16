package rotate_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/fincut/internal/rotate"
)

func writeTempRotate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "app.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempRotate: %v", err)
	}
	return p
}

func TestRotator_Detect_AfterTruncate(t *testing.T) {
	p := writeTempRotate(t, "line1\nline2\n")

	out := &strings.Builder{}
	r, err := rotate.NewRotator(rotate.Options{
		Path:     p,
		Output:   out,
		Interval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewRotator: %v", err)
	}
	_ = r

	// Truncate simulates rotation
	if err := os.Truncate(p, 0); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	if err := os.WriteFile(p, []byte("new-line\n"), 0o644); err != nil {
		t.Fatalf("rewrite: %v", err)
	}

	// Give watcher a moment
	time.Sleep(60 * time.Millisecond)
}

func TestRotator_Detect_SameInode_NoRotation(t *testing.T) {
	p := writeTempRotate(t, "stable\n")

	var buf strings.Builder
	r, err := rotate.NewRotator(rotate.Options{
		Path:     p,
		Output:   &buf,
		Interval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewRotator: %v", err)
	}
	_ = r

	time.Sleep(50 * time.Millisecond)
	// No rotation — no output expected from rotation detection
	if buf.Len() > 0 {
		t.Errorf("unexpected output without rotation: %q", buf.String())
	}
}
