package prefix

import "errors"

func validateOptions(o *Options) error {
	if !o.LineNumbers && o.Text == "" {
		return errors.New("prefix: either Text or LineNumbers must be set")
	}
	if o.LineNumbers && o.Text != "" {
		return errors.New("prefix: Text and LineNumbers are mutually exclusive")
	}
	if o.Width < 0 {
		return errors.New("prefix: Width must be non-negative")
	}
	if o.Separator == "" {
		o.Separator = ": "
	}
	return nil
}
