package split

import "errors"

func validateOptions(opts Options) error {
	if opts.Pattern == "" {
		return errors.New("split: pattern must not be empty")
	}
	return nil
}
