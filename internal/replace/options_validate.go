package replace

import "errors"

// validateOptions returns an error if opts is not usable.
func validateOptions(opts Options) error {
	if len(opts.Patterns) == 0 {
		return errors.New("replace: at least one pattern=replacement pair is required")
	}
	for _, p := range opts.Patterns {
		hasEq := false
		for _, c := range p {
			if c == '=' {
				hasEq = true
				break
			}
		}
		if !hasEq {
			return errors.New("replace: each entry must be in 'pattern=replacement' form")
		}
	}
	return nil
}
