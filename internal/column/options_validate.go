package column

import "errors"

func validateOptions(opts Options) error {
	if opts.Delimiter == "" {
		return errors.New("column: delimiter must not be empty")
	}
	if len(opts.Fields) == 0 {
		return errors.New("column: at least one field index is required")
	}
	for _, f := range opts.Fields {
		if f < 1 {
			return errors.New("column: field indices must be >= 1")
		}
	}
	return nil
}
