// Package field provides utilities for extracting named fields from
// structured log lines using delimiter-based or regex-based parsing.
package field

import (
	"fmt"
	"regexp"
	"strings"
)

// Options configures the field extractor.
type Options struct {
	// Delimiter splits each line into fields (e.g. " ", "\t", "|").
	// Ignored when Pattern is set.
	Delimiter string

	// Pattern is a named-capture regex used to extract fields.
	// When set, Delimiter is ignored.
	Pattern string

	// Fields is the ordered list of field names to emit.
	// For delimiter mode these are positional labels (index 0 = first token).
	Fields []string
}

// Extractor pulls named fields out of log lines.
type Extractor struct {
	opts Options
	re   *regexp.Regexp
}

// New returns a validated Extractor or an error.
func New(opts Options) (*Extractor, error) {
	if opts.Pattern == "" && opts.Delimiter == "" {
		return nil, fmt.Errorf("field: either Pattern or Delimiter must be set")
	}
	if len(opts.Fields) == 0 {
		return nil, fmt.Errorf("field: Fields must not be empty")
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("field: invalid pattern: %w", err)
		}
		// Ensure every requested field has a named capture group.
		names := make(map[string]bool)
		for _, n := range re.SubexpNames() {
			names[n] = true
		}
		for _, f := range opts.Fields {
			if !names[f] {
				return nil, fmt.Errorf("field: pattern has no capture group %q", f)
			}
		}
	}

	return &Extractor{opts: opts, re: re}, nil
}

// Extract returns a map of field name → value for the given line.
// Missing fields are represented as empty strings.
func (e *Extractor) Extract(line string) map[string]string {
	out := make(map[string]string, len(e.opts.Fields))
	for _, f := range e.opts.Fields {
		out[f] = ""
	}

	if e.re != nil {
		m := e.re.FindStringSubmatch(line)
		if m == nil {
			return out
		}
		for i, name := range e.re.SubexpNames() {
			if name != "" {
				out[name] = m[i]
			}
		}
		return out
	}

	// Delimiter mode.
	tokens := strings.Split(line, e.opts.Delimiter)
	for i, f := range e.opts.Fields {
		if i < len(tokens) {
			out[f] = tokens[i]
		}
	}
	return out
}
