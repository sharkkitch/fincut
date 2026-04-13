package burst

import (
	"errors"
	"time"
)

// Options configures burst detection behaviour.
type Options struct {
	// TimestampPattern is a regex used to locate the timestamp in each line.
	TimestampPattern string
	// TimestampLayout is the time.Parse layout for the matched timestamp.
	TimestampLayout string
	// Window is the sliding duration over which rate is measured.
	Window time.Duration
	// Threshold is the minimum lines-per-second rate that constitutes a burst.
	Threshold float64
}

func validateOptions(o *Options) error {
	if o.TimestampPattern == "" {
		return errors.New("burst: TimestampPattern must not be empty")
	}
	if o.TimestampLayout == "" {
		return errors.New("burst: TimestampLayout must not be empty")
	}
	if o.Window <= 0 {
		return errors.New("burst: Window must be a positive duration")
	}
	if o.Threshold <= 0 {
		return errors.New("burst: Threshold must be greater than zero")
	}
	return nil
}
