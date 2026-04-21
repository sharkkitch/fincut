package echo

import "errors"

func validateOptions(opts Options) error {
	if opts.Writer == nil {
		return errors.New("echo: Writer must not be nil")
	}
	return nil
}
