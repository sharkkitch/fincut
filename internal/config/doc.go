// Package config provides YAML-based configuration loading and validation
// for fincut. It defines the Config struct that controls input file selection,
// byte-range offsets, regex filter pipeline patterns, output formatting, and
// statistics reporting.
//
// Usage:
//
//	cfg, err := config.Load(".fincut.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// A DefaultConfig is available for programmatic construction, and Merge
// allows CLI flags to override file-based settings without mutating the
// original config.
package config
