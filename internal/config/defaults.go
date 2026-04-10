package config

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Input: InputConfig{
			Files:       []string{},
			StartOffset: 0,
			EndOffset:   0,
		},
		Filter: FilterConfig{
			Patterns: []string{},
		},
		Output: OutputConfig{
			Format: "plain",
		},
		Stats: StatsConfig{
			Enabled: false,
			Format:  "plain",
		},
	}
}

// Merge overlays non-zero values from override onto base, returning a new Config.
// Fields in override take precedence when they carry a non-zero value.
func Merge(base, override *Config) *Config {
	out := *base

	if len(override.Input.Files) > 0 {
		out.Input.Files = override.Input.Files
	}
	if override.Input.StartOffset != 0 {
		out.Input.StartOffset = override.Input.StartOffset
	}
	if override.Input.EndOffset != 0 {
		out.Input.EndOffset = override.Input.EndOffset
	}
	if len(override.Filter.Patterns) > 0 {
		out.Filter.Patterns = override.Filter.Patterns
	}
	if override.Output.Format != "" {
		out.Output.Format = override.Output.Format
	}
	if override.Stats.Enabled {
		out.Stats.Enabled = true
	}
	if override.Stats.Format != "" {
		out.Stats.Format = override.Stats.Format
	}

	return &out
}
