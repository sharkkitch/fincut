package normalize

import "errors"

// validateOptions ensures that at least one normalization option is enabled.
func validateOptions(opts Options) error {
	if !opts.TrimSpace && !opts.CollapseSpaces && !opts.Lowercase && !opts.StripControl {
		return errors.New("normalize: at least one normalization option must be enabled")
	}
	return nil
}
