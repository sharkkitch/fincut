package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level fincut configuration.
type Config struct {
	Input   InputConfig   `yaml:"input"`
	Filter  FilterConfig  `yaml:"filter"`
	Output  OutputConfig  `yaml:"output"`
	Stats   StatsConfig   `yaml:"stats"`
}

// InputConfig controls how log files are read.
type InputConfig struct {
	Files       []string `yaml:"files"`
	StartOffset int64    `yaml:"start_offset"`
	EndOffset   int64    `yaml:"end_offset"`
}

// FilterConfig defines the regex pipeline stages.
type FilterConfig struct {
	Patterns []string `yaml:"patterns"`
}

// OutputConfig controls output formatting.
type OutputConfig struct {
	Format string `yaml:"format"` // plain, json, color
}

// StatsConfig controls statistics reporting.
type StatsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"` // plain, json
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

// validate checks that required fields are set and values are in range.
func (c *Config) validate() error {
	validFormats := map[string]bool{"plain": true, "json": true, "color": true}
	if c.Output.Format != "" && !validFormats[c.Output.Format] {
		return fmt.Errorf("output.format %q is not one of plain, json, color", c.Output.Format)
	}
	if c.Stats.Format != "" {
		if c.Stats.Format != "plain" && c.Stats.Format != "json" {
			return fmt.Errorf("stats.format %q is not one of plain, json", c.Stats.Format)
		}
	}
	if c.Input.StartOffset < 0 {
		return fmt.Errorf("input.start_offset must be >= 0")
	}
	if c.Input.EndOffset < 0 {
		return fmt.Errorf("input.end_offset must be >= 0")
	}
	if c.Input.EndOffset > 0 && c.Input.EndOffset <= c.Input.StartOffset {
		return fmt.Errorf("input.end_offset must be greater than input.start_offset")
	}
	return nil
}
