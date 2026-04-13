package ratelimit

import (
	"errors"
	"time"
)

func validateOptions(opts Options) error {
	if opts.MaxLines <= 0 {
		return errors.New("ratelimit: MaxLines must be greater than zero")
	}
	if opts.Window <= 0 {
		return errors.New("ratelimit: Window must be a positive duration")
	}
	if opts.Window < time.Millisecond {
		return errors.New("ratelimit: Window must be at least 1ms")
	}
	return nil
}
