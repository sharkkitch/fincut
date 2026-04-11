package annotate

import (
	"errors"
	"time"
)

// Options controls the behaviour of the Annotator.
type Options struct {
	// LineNumbers prepends a 1-based line number to each line.
	LineNumbers bool

	// Timestamp prepends the current UTC time to each line.
	Timestamp bool

	// TimestampFormat is the Go time format string used when Timestamp is true.
	// Defaults to time.RFC3339 if empty.
	TimestampFormat string

	// Prefix is a static string prepended before the line content.
	Prefix string

	// Separator is placed between annotation fields and the line body.
	// Defaults to " | " if empty.
	Separator string
}

func validateOptions(o Options) error {
	if !o.LineNumbers && !o.Timestamp && o.Prefix == "" {
		return errors.New("annotate: at least one of LineNumbers, Timestamp, or Prefix must be set")
	}

	if o.Timestamp && o.TimestampFormat != "" {
		// Validate the format by attempting a round-trip.
		formatted := time.Now().UTC().Format(o.TimestampFormat)
		if formatted == "" {
			return errors.New("annotate: TimestampFormat produced empty output")
		}
	}

	return nil
}

// defaults fills in zero-value fields with sensible defaults.
func defaults(o *Options) {
	if o.TimestampFormat == "" {
		o.TimestampFormat = time.RFC3339
	}
	if o.Separator == "" {
		o.Separator = " | "
	}
}
