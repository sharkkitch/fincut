package config_test

import (
	"testing"

	"github.com/yourorg/fincut/internal/config"
)

// TestMerge_CLIOverridesFile simulates the common pattern of loading a file-based
// config and then applying CLI flag overrides on top.
func TestMerge_CLIOverridesFile(t *testing.T) {
	fileCfg := writeTemp(t, `
input:
  files: ["server.log"]
output:
  format: plain
stats:
  enabled: false
  format: plain
`)

	base, err := config.Load(fileCfg)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	// Simulate CLI flags: user passes --format=color --stats
	cliOverride := &config.Config{
		Output: config.OutputConfig{Format: "color"},
		Stats:  config.StatsConfig{Enabled: true, Format: "json"},
	}

	merged := config.Merge(base, cliOverride)

	if merged.Output.Format != "color" {
		t.Errorf("expected color output, got %q", merged.Output.Format)
	}
	if !merged.Stats.Enabled {
		t.Error("expected stats enabled after merge")
	}
	if merged.Stats.Format != "json" {
		t.Errorf("expected json stats format, got %q", merged.Stats.Format)
	}
	// Original file-based input should be preserved
	if len(merged.Input.Files) != 1 || merged.Input.Files[0] != "server.log" {
		t.Errorf("expected input files preserved, got %v", merged.Input.Files)
	}
}

// TestMerge_EmptyOverridePreservesBase ensures that an empty override struct
// does not wipe out base values.
func TestMerge_EmptyOverridePreservesBase(t *testing.T) {
	base := config.DefaultConfig()
	base.Output.Format = "json"
	base.Stats.Enabled = true

	merged := config.Merge(base, &config.Config{})

	if merged.Output.Format != "json" {
		t.Errorf("expected json preserved, got %q", merged.Output.Format)
	}
	if !merged.Stats.Enabled {
		t.Error("expected stats.enabled preserved")
	}
}
