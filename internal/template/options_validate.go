package template

import "errors"

func validateOptions(opts Options) error {
	if opts.Template == "" {
		return errors.New("template: Template must not be empty")
	}
	return nil
}
