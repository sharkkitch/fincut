package sparse

import "fmt"

func validateOptions(opts Options) error {
	if opts.Every < 1 {
		return fmt.Errorf("sparse: Every must be >= 1, got %d", opts.Every)
	}
	if opts.Offset < 0 {
		return fmt.Errorf("sparse: Offset must be >= 0, got %d", opts.Offset)
	}
	if opts.Offset >= opts.Every {
		return fmt.Errorf("sparse: Offset (%d) must be less than Every (%d)", opts.Offset, opts.Every)
	}
	return nil
}
