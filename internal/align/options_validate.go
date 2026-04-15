package align

import "errors"

func validateOptions(opts *Options) error {
	if opts.Delimiter == "" {
		return errors.New("align: delimiter must not be empty")
	}
	if opts.Padding < 0 {
		return errors.New("align: padding must not be negative")
	}
	if opts.Padding == 0 {
		opts.Padding = 1
	}
	return nil
}
