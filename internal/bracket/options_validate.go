package bracket

import "errors"

func validateOptions(opts Options) error {
	if !opts.WrapAll && opts.Pattern == "" {
		return errors.New("bracket: pattern must be set when WrapAll is false")
	}
	if opts.Open == "" && opts.Close == "" {
		return errors.New("bracket: at least one of Open or Close must be non-empty")
	}
	return nil
}
