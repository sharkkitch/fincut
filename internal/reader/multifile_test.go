package reader_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/fincut/internal/reader"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestNewMultiFileReader_EmptyPaths(t *testing.T) {
	_, err := reader.NewMultiFileReader([]string{})
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestNewMultiFileReader_MissingFile(t *testing.T) {
	_, err := reader.NewMultiFileReader([]string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestMultiFileReader_Lines_SingleFile(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "a.log", "line1\nline2\nline3\n")

	mfr, err := reader.NewMultiFileReader([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines, err := mfr.Lines()
	if err != nil {
		t.Fatalf("Lines() error: %v", err)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, l := range lines {
		if l != expected[i] {
			t.Errorf("line[%d]: want %q, got %q", i, expected[i], l)
		}
	}
}

func TestMultiFileReader_Lines_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempFile(t, dir, "first.log", "alpha\nbeta\n")
	p2 := writeTempFile(t, dir, "second.log", "gamma\ndelta\n")

	mfr, err := reader.NewMultiFileReader([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines, err := mfr.Lines()
	if err != nil {
		t.Fatalf("Lines() error: %v", err)
	}

	expected := []string{"alpha", "beta", "gamma", "delta"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i, l := range lines {
		if l != expected[i] {
			t.Errorf("line[%d]: want %q, got %q", i, expected[i], l)
		}
	}
}

func TestMultiFileReader_Lines_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "empty.log", "")

	mfr, err := reader.NewMultiFileReader([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines, err := mfr.Lines()
	if err != nil {
		t.Fatalf("Lines() error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d: %v", len(lines), strings.Join(lines, ","))
	}
}
