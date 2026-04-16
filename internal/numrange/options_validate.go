package numrange

import "fmt"

func validateOptions(opts Options) error {
	if opts.Pattern == "" {
		return fmt.Errorf("numrange: pattern must not be empty")
	}
	if opts.Min != nil && opts.Max != nil && *opts.Min > *opts.Max {
		return fmt.Errorf("numrange: min (%.4g) must not exceed max (%.4g)", *opts.Min, *opts.Max)
	}
	return nil
}
