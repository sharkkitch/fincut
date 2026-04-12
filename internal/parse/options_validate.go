package parse

import "fmt"

func validateOptions(opts Options) error {
	switch opts.Format {
	case FormatJSON:
		// no additional fields required
	case FormatRegex:
		if strings.TrimSpace(opts.Pattern) == "" {
			return fmt.Errorf("parse: pattern is required for regex format")
		}
	case FormatDelim:
		if opts.Delimiter == "" {
			return fmt.Errorf("parse: delimiter is required for delim format")
		}
		if len(opts.Fields) == 0 {
			return fmt.Errorf("parse: at least one field name is required for delim format")
		}
	case "":
		return fmt.Errorf("parse: format must be one of: json, regex, delim")
	default:
		return fmt.Errorf("parse: unknown format %q; must be one of: json, regex, delim", opts.Format)
	}
	return nil
}
