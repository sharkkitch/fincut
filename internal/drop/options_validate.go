package drop

import "errors"

func validateOptions(opts Options) error {
	if len(opts.Patterns) == 0 {
		return errors.New("drop: at least one pattern is required")
	}
	return nil
}
