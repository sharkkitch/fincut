// Package tokenize splits log lines into tokens using a delimiter or regex,
// enabling per-token analysis and transformation.
package tokenize

import (
	"fmt"
	"regexp"
	"strings"
)

// Options configures the Tokenizer.
type Options struct {
	// Delimiter splits each line into tokens. Mutually exclusive with Pattern.
	Delimiter string
	// Pattern is a regex whose submatches become tokens. Mutually exclusive with Delimiter.
	Pattern string
	// MinTokens discards lines that produce fewer tokens than this value. 0 = no minimum.
	MinTokens int
	// Join is the string used to reassemble tokens into the output line.
	Join string
}

// Tokenizer splits lines into tokens and reassembles them.
type Tokenizer struct {
	opts Options
	re   *regexp.Regexp
}

// New creates a Tokenizer from opts, returning an error if the configuration is invalid.
func New(opts Options) (*Tokenizer, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	t := &Tokenizer{opts: opts}
	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("tokenize: invalid pattern: %w", err)
		}
		if re.NumSubexp() == 0 {
			return nil, fmt.Errorf("tokenize: pattern must contain at least one capture group")
		}
		t.re = re
	}
	return t, nil
}

// Apply tokenizes each line and returns the reassembled results.
// Lines producing fewer tokens than MinTokens are dropped.
func (t *Tokenizer) Apply(lines []string) []string {
	join := t.opts.Join
	if join == "" {
		join = " "
	}
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		tokens := t.tokenize(line)
		if t.opts.MinTokens > 0 && len(tokens) < t.opts.MinTokens {
			continue
		}
		out = append(out, strings.Join(tokens, join))
	}
	return out
}

func (t *Tokenizer) tokenize(line string) []string {
	if t.re != nil {
		m := t.re.FindStringSubmatch(line)
		if m == nil {
			return nil
		}
		return m[1:]
	}
	return strings.Split(line, t.opts.Delimiter)
}
