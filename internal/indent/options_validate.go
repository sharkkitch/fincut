package indent

import "errors"

func validateOptions(opts *Options) error {
	if opts.Depth == 0 && !opts.StripExisting {
		return errors.New("indent: depth is zero and strip_existing is false — no transformation would occur")
	}
	if opts.Unit == "" {
		opts.Unit = "  "
	}
	return nil
}
