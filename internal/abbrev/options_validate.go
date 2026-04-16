package abbrev

import "errors"

func validateOptions(opts *Options) error {
	if opts.MaxTokenLen <= 0 {
		return errors.New("abbrev: MaxTokenLen must be greater than zero")
	}
	if opts.Delimiter == "" {
		opts.Delimiter = " "
	}
	if opts.Ellipsis == "" {
		opts.Ellipsis = "…"
	}
	return nil
}
