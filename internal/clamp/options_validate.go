package clamp

import "fmt"

func validateOptions(opts Options) error {
	if opts.Pattern == "" {
		return fmt.Errorf("clamp: pattern must not be empty")
	}
	if opts.Min != nil && opts.Max != nil && *opts.Min > *opts.Max {
		return fmt.Errorf("clamp: min (%g) must not exceed max (%g)", *opts.Min, *opts.Max)
	}
	if opts.Min == nil && opts.Max == nil {
		return fmt.Errorf("clamp: at least one of Min or Max must be set")
	}
	return nil
}
