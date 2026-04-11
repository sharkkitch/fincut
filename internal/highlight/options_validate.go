package highlight

import (
	"fmt"
	"regexp"
)

// ValidatePatterns checks that every pattern in the slice compiles as a valid
// regular expression. Returns the first error encountered, or nil.
func ValidatePatterns(patterns []string) error {
	if len(patterns) == 0 {
		return fmt.Errorf("highlight: patterns must not be empty")
	}
	for _, p := range patterns {
		if p == "" {
			return fmt.Errorf("highlight: pattern must not be an empty string")
		}
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("highlight: invalid pattern %q: %w", p, err)
		}
	}
	return nil
}

// MaxPatterns is the upper bound on the number of simultaneous highlight
// patterns to keep output legible.
const MaxPatterns = 8

// validateOptions checks Options before a Highlighter is constructed.
func validateOptions(opts Options) error {
	if err := ValidatePatterns(opts.Patterns); err != nil {
		return err
	}
	if len(opts.Patterns) > MaxPatterns {
		return fmt.Errorf("highlight: too many patterns (%d); maximum is %d",
			len(opts.Patterns), MaxPatterns)
	}
	return nil
}
