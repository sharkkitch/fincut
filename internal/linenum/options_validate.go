package linenum

import "fmt"

func validateOptions(opts Options) error {
	if len(opts.Ranges) == 0 {
		return fmt.Errorf("linenum: at least one range is required")
	}
	for i, r := range opts.Ranges {
		start, end := r[0], r[1]
		if start < 1 {
			return fmt.Errorf("linenum: range[%d] start must be >= 1, got %d", i, start)
		}
		if end != 0 && end < start {
			return fmt.Errorf("linenum: range[%d] end (%d) must be >= start (%d) or 0 for EOF", i, end, start)
		}
	}
	return nil
}
