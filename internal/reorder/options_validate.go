package reorder

import "fmt"

func validateOptions(opts Options) error {
	if opts.Field < 0 {
		return fmt.Errorf("reorder: field index must be >= 0, got %d", opts.Field)
	}
	if opts.Field > 0 && opts.Delimiter == "" {
		return fmt.Errorf("reorder: delimiter must be set when field > 0")
	}
	if !opts.Reverse && opts.Field == 0 {
		return fmt.Errorf("reorder: at least one of reverse or field must be set")
	}
	return nil
}
