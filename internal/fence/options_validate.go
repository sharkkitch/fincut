package fence

import "errors"

func validateOptions(opts Options) error {
	if opts.OpenPattern == "" {
		return errors.New("fence: open pattern must not be empty")
	}
	if opts.ClosePattern == "" {
		return errors.New("fence: close pattern must not be empty")
	}
	return nil
}
