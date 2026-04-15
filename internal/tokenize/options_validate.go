package tokenize

import "fmt"

func validateOptions(opts Options) error {
	if opts.Delimiter == "" && opts.Pattern == "" {
		return fmt.Errorf("tokenize: one of Delimiter or Pattern must be set")
	}
	if opts.Delimiter != "" && opts.Pattern != "" {
		return fmt.Errorf("tokenize: Delimiter and Pattern are mutually exclusive")
	}
	if opts.MinTokens < 0 {
		return fmt.Errorf("tokenize: MinTokens must be non-negative, got %d", opts.MinTokens)
	}
	return nil
}
