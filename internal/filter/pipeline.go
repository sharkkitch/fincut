package filter

import (
	"fmt"
	"regexp"
)

// Stage represents a single regex-based filter stage in the pipeline.
type Stage struct {
	Pattern *regexp.Regexp
	Invert  bool // if true, keep lines that do NOT match
}

// Pipeline is an ordered sequence of filter stages applied to log lines.
type Pipeline struct {
	Stages []Stage
}

// NewPipeline constructs a Pipeline from a slice of pattern strings.
// Patterns prefixed with '!' are treated as inverted (exclude) filters.
func NewPipeline(patterns []string) (*Pipeline, error) {
	p := &Pipeline{}
	for _, raw := range patterns {
		invert := false
		if len(raw) > 0 && raw[0] == '!' {
			invert = true
			raw = raw[1:]
		}
		re, err := regexp.Compile(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", raw, err)
		}
		p.Stages = append(p.Stages, Stage{Pattern: re, Invert: invert})
	}
	return p, nil
}

// Match returns true if the given line passes all stages in the pipeline.
func (p *Pipeline) Match(line string) bool {
	for _, s := range p.Stages {
		matched := s.Pattern.MatchString(line)
		if s.Invert {
			if matched {
				return false
			}
		} else {
			if !matched {
				return false
			}
		}
	}
	return true
}

// Apply filters a slice of lines through the pipeline and returns matching lines.
func (p *Pipeline) Apply(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		if p.Match(l) {
			result = append(result, l)
		}
	}
	return result
}
