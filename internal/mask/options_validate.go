package mask

import (
	"fmt"
	"regexp"
)

// validateOptions checks that the provided Options are well-formed.
func validateOptions(opts Options) error {
	if len(opts.Patterns) == 0 {
		return fmt.Errorf("mask: at least one pattern is required")
	}
	for _, p := range opts.Patterns {
		if p == "" {
			return fmt.Errorf("mask: pattern must not be empty")
		}
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("mask: invalid pattern %q: %w", p, err)
		}
	}
	if len(opts.Replacement) > 256 {
		return fmt.Errorf("mask: replacement string exceeds maximum length of 256 characters")
	}
	return nil
}
