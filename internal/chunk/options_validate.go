package chunk

import "errors"

func validateOptions(o *Options) error {
	if o.Size > 0 && o.Delimiter != "" {
		return errors.New("chunk: Size and Delimiter are mutually exclusive")
	}
	if o.Size == 0 && o.Delimiter == "" {
		return errors.New("chunk: one of Size or Delimiter must be set")
	}
	if o.Size < 0 {
		return errors.New("chunk: Size must be positive")
	}
	if o.LabelPrefix == "" {
		o.LabelPrefix = "chunk"
	}
	return nil
}
