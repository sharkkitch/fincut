package summarize

import "errors"

func validateOptions(opts Options) error {
	if opts.TopN < 0 {
		return errors.New("summarize: TopN must be non-negative")
	}
	if opts.MinCount < 0 {
		return errors.New("summarize: MinCount must be non-negative")
	}
	return nil
}
