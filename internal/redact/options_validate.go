package redact

import "errors"

func validateOptions(o *Options) error {
	if len(o.Patterns) == 0 {
		return errors.New("redact: at least one pattern is required")
	}
	if o.Replacement == "" {
		o.Replacement = "[REDACTED]"
	}
	return nil
}
