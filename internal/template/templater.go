package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	text "text/template"
)

// Templater applies a Go text/template to each input line, exposing named
// capture groups from an optional regex as template variables.
type Templater struct {
	templ   *text.Template
	pattern *regexp.Regexp
}

// Options configures a Templater.
type Options struct {
	// Template is the Go text/template string applied to each line.
	Template string
	// Pattern is an optional regex with named capture groups. When set,
	// capture group values are exposed as template variables. Lines that
	// do not match are passed through unchanged.
	Pattern string
}

// New creates a Templater from the provided options.
func New(opts Options) (*Templater, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	t, err := text.New("line").Option("missingkey=zero").Parse(opts.Template)
	if err != nil {
		return nil, fmt.Errorf("template: parse error: %w", err)
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("template: invalid pattern: %w", err)
		}
	}

	return &Templater{templ: t, pattern: re}, nil
}

// Apply executes the template against each line and returns the rendered lines.
func (t *Templater) Apply(lines []string) ([]string, error) {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		rendered, err := t.render(line)
		if err != nil {
			return nil, err
		}
		out = append(out, rendered)
	}
	return out, nil
}

func (t *Templater) render(line string) (string, error) {
	data := map[string]string{"Line": line}

	if t.pattern != nil {
		m := t.pattern.FindStringSubmatch(line)
		if m == nil {
			return line, nil
		}
		for i, name := range t.pattern.SubexpNames() {
			if name != "" {
				data[name] = m[i]
			}
		}
	}

	var buf bytes.Buffer
	if err := t.templ.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template: execute error: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}
