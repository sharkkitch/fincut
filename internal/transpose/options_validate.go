package transpose

import "errors"

func validateOptions(opts Options) error {
	if opts.Delimiter == "" {
		return errors.New("transpose: delimiter must not be empty")
	}
	return nil
}
