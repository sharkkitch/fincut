package watch

import (
	"fmt"
	"os"
	"time"
)

const (
	// MinPollInterval is the smallest allowed polling interval.
	MinPollInterval = 10 * time.Millisecond
	// MaxPollInterval is the largest allowed polling interval.
	MaxPollInterval = 60 * time.Second
)

// Validate checks that Options are coherent before constructing a Watcher.
func (o Options) Validate() error {
	if o.Path == "" {
		return fmt.Errorf("watch: path must not be empty")
	}
	if o.Output == nil {
		return fmt.Errorf("watch: output writer must not be nil")
	}
	if o.PollInterval != 0 && o.PollInterval < MinPollInterval {
		return fmt.Errorf("watch: poll interval %v is below minimum %v", o.PollInterval, MinPollInterval)
	}
	if o.PollInterval > MaxPollInterval {
		return fmt.Errorf("watch: poll interval %v exceeds maximum %v", o.PollInterval, MaxPollInterval)
	}
	return nil
}

// fileExists returns true when path refers to a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
