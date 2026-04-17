package stripe

import "fmt"

func validateOptions(opts Options) error {
	if opts.Every < 2 {
		return fmt.Errorf("stripe: every must be >= 2, got %d", opts.Every)
	}
	if opts.Offset < 0 {
		return fmt.Errorf("stripe: offset must be >= 0, got %d", opts.Offset)
	}
	if opts.Offset >= opts.Every {
		return fmt.Errorf("stripe: offset %d must be less than every %d", opts.Offset, opts.Every)
	}
	return nil
}
