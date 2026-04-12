package join

import "fmt"

// Options configures the Joiner.
type Options struct {
	// GroupSize is the number of consecutive lines to join into one.
	// A value of 0 means join all lines into a single output line.
	GroupSize int

	// Separator is the string placed between joined lines.
	// Defaults to a single space if empty.
	Separator string
}

func validateOptions(opts *Options) error {
	if opts.GroupSize < 0 {
		return fmt.Errorf("join: GroupSize must be >= 0, got %d", opts.GroupSize)
	}
	if opts.Separator == "" {
		opts.Separator = " "
	}
	return nil
}
