package timestamp

import (
	"errors"
	"time"
)

func validateOptions(o *Options) error {
	if o.Prepend && o.Append {
		return errors.New("timestamp: Prepend and Append are mutually exclusive")
	}
	if !o.Prepend && !o.Append {
		return errors.New("timestamp: one of Prepend or Append must be set")
	}
	if o.Format == "" {
		o.Format = time.RFC3339
	}
	if o.Sep == "" {
		o.Sep = " "
	}
	return nil
}
