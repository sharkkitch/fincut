package between

import "errors"

func validateOptions(opts Options) error {
	if opts.StartPattern == "" {
		return errors.New("between: start pattern must not be empty")
	}
	if opts.EndPattern == "" {
		return errors.New("between: end pattern must not be empty")
	}
	return nil
}
