package fold

import "errors"

func validateOptions(o *Options) error {
	if o.ContinuationPattern == "" {
		return errors.New("fold: ContinuationPattern must not be empty")
	}
	if o.Separator == "" {
		o.Separator = " "
	}
	return nil
}
