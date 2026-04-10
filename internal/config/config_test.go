package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/fincut/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTemp(t, `
input:
  files: ["a.log", "b.log"]
  start_offset: 0
  end_offset: 1024
filter:
  patterns: ["ERROR", "WARN"]
output:
  format: json
stats:
  enabled: true
  format: plain
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Input.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Input.Files))
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected format json, got %q", cfg.Output.Format)
	}
	if !cfg.Stats.Enabled {
		t.Error("expected stats.enabled to be true")
	}
}

func TestLoad_InvalidOutputFormat(t *testing.T) {
	p := writeTemp(t, `output:\n  format: xml\n`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid output format")
	}
}

func TestLoad_InvalidOffsets(t *testing.T) {
	p := writeTemp(t, `
input:
  start_offset: 500
  end_offset: 100
`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error when end_offset <= start_offset")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Output.Format != "plain" {
		t.Errorf("expected default format plain, got %q", cfg.Output.Format)
	}
	if cfg.Stats.Enabled {
		t.Error("expected stats disabled by default")
	}
}

func TestMerge_OverrideFormat(t *testing.T) {
	base := config.DefaultConfig()
	override := &config.Config{
		Output: config.OutputConfig{Format: "color"},
	}
	merged := config.Merge(base, override)
	if merged.Output.Format != "color" {
		t.Errorf("expected color, got %q", merged.Output.Format)
	}
	if merged.Stats.Format != base.Stats.Format {
		t.Errorf("expected stats format preserved, got %q", merged.Stats.Format)
	}
}
