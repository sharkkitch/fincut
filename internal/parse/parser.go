// Package parse provides structured log line parsing into key-value fields.
// It supports JSON log lines as well as delimiter-separated and regex-named-group formats.
package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Format specifies how log lines should be parsed.
type Format string

const (
	FormatJSON    Format = "json"
	FormatRegex   Format = "regex"
	FormatDelim   Format = "delim"
)

// Options configures the Parser.
type Options struct {
	Format    Format
	Pattern   string // used with FormatRegex; must contain named capture groups
	Delimiter string // used with FormatDelim
	Fields    []string // field names for FormatDelim columns
}

// Parser parses log lines into field maps.
type Parser struct {
	opts Options
	re   *regexp.Regexp
}

// New constructs a Parser from the given Options, returning an error if the
// configuration is invalid.
func New(opts Options) (*Parser, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	p := &Parser{opts: opts}
	if opts.Format == FormatRegex {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("parse: invalid pattern: %w", err)
		}
		if len(re.SubexpNames()) < 2 {
			return nil, fmt.Errorf("parse: pattern must contain at least one named capture group")
		}
		p.re = re
	}
	return p, nil
}

// Parse parses a single log line and returns a map of field names to values.
// Returns an error if the line cannot be parsed according to the configured format.
func (p *Parser) Parse(line string) (map[string]string, error) {
	switch p.opts.Format {
	case FormatJSON:
		return parseJSON(line)
	case FormatRegex:
		return p.parseRegex(line)
	case FormatDelim:
		return p.parseDelim(line)
	default:
		return nil, fmt.Errorf("parse: unknown format %q", p.opts.Format)
	}
}

func parseJSON(line string) (map[string]string, error) {
	raw := map[string]interface{}{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("parse: json: %w", err)
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

func (p *Parser) parseRegex(line string) (map[string]string, error) {
	match := p.re.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("parse: line did not match pattern")
	}
	out := make(map[string]string)
	for i, name := range p.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}
	return out, nil
}

func (p *Parser) parseDelim(line string) (map[string]string, error) {
	parts := strings.Split(line, p.opts.Delimiter)
	out := make(map[string]string, len(p.opts.Fields))
	for i, field := range p.opts.Fields {
		if i < len(parts) {
			out[field] = parts[i]
		} else {
			out[field] = ""
		}
	}
	return out, nil
}
