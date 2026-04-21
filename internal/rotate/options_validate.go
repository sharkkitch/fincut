package rotate

import (
	"errors"
)

func validateOptions(opts Options) error {
	if opts.Path == "" {
		return errors.New("rotate: path must not be empty")
	}
	if opts.Out == nil {
		return errors.New("rotate: output writer must not be nil")
	}
	return nil
}

func fileExists(path string) bool {
	_, err := statFile(path)
	return err == nil
}
