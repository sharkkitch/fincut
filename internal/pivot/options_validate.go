package pivot

import "fmt"

// Options configures a Pivotter.
type Options struct {
	// KeyField is the zero-based column index used as the grouping key.
	KeyField int
	// ValueField is the zero-based column index to aggregate. Use -1 to ignore.
	ValueField int
	// Delimiter separates fields in each line.
	Delimiter string
	// Aggregator is one of: count, sum, values.
	Aggregator string
}

func validateOptions(opts Options) error {
	if opts.KeyField < 0 {
		return fmt.Errorf("pivot: key_field must be >= 0, got %d", opts.KeyField)
	}
	if opts.ValueField < -1 {
		return fmt.Errorf("pivot: value_field must be >= -1, got %d", opts.ValueField)
	}
	if opts.Delimiter == "" {
		return fmt.Errorf("pivot: delimiter must not be empty")
	}
	valid := map[string]bool{"count": true, "sum": true, "values": true, "": true}
	if !valid[opts.Aggregator] {
		return fmt.Errorf("pivot: unknown aggregator %q; want count, sum, or values", opts.Aggregator)
	}
	return nil
}
