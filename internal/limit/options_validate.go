package limit

import "errors"

func validateOptions(opts Options) error {
	if opts.MaxLines < 0 {
		return errors.New("limit: MaxLines must not be negative")
	}
	if opts.MaxBytes < 0 {
		return errors.New("limit: MaxBytes must not be negative")
	}
	if opts.MaxLines == 0 && opts.MaxBytes == 0 {
		return errors.New("limit: at least one of MaxLines or MaxBytes must be set")
	}
	return nil
}
