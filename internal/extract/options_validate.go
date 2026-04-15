package extract

import "fmt"

func validateOptions(o *Options) error {
	if o.Pattern == "" {
		return fmt.Errorf("extract: pattern must not be empty")
	}
	if o.Group < 0 {
		return fmt.Errorf("extract: group must be >= 0, got %d", o.Group)
	}
	if o.Group == 0 {
		o.Group = 1
	}
	return nil
}
