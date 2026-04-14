package pad

import "errors"

func validateOptions(opts *Options) error {
	if opts.Width <= 0 {
		return errors.New("pad: width must be greater than zero")
	}
	if opts.Fill == 0 {
		opts.Fill = ' '
	}
	return nil
}
