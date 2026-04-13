// Package chunk splits a stream of lines into fixed-size or delimiter-bounded
// chunks, emitting each chunk as a named group for downstream processing.
package chunk

import (
	"fmt"
	"regexp"
)

// Chunk holds a labeled group of lines produced by the Chunker.
type Chunk struct {
	Label string
	Lines []string
}

// Options configures the Chunker.
type Options struct {
	// Size is the maximum number of lines per chunk (mutually exclusive with Delimiter).
	Size int
	// Delimiter is a regex; a matching line starts a new chunk (mutually exclusive with Size).
	Delimiter string
	// LabelPrefix is prepended to the chunk index in the label (default "chunk").
	LabelPrefix string
}

// Chunker splits lines into chunks.
type Chunker struct {
	opts Options
	delim *regexp.Regexp
}

// New constructs a Chunker from opts, returning an error if the configuration
// is invalid.
func New(opts Options) (*Chunker, error) {
	if err := validateOptions(&opts); err != nil {
		return nil, err
	}
	var re *regexp.Regexp
	if opts.Delimiter != "" {
		var err error
		re, err = regexp.Compile(opts.Delimiter)
		if err != nil {
			return nil, fmt.Errorf("chunk: invalid delimiter pattern: %w", err)
		}
	}
	return &Chunker{opts: opts, delim: re}, nil
}

// Apply splits lines into chunks according to the configured strategy.
func (c *Chunker) Apply(lines []string) []Chunk {
	if len(lines) == 0 {
		return nil
	}
	if c.delim != nil {
		return c.splitByDelimiter(lines)
	}
	return c.splitBySize(lines)
}

func (c *Chunker) splitBySize(lines []string) []Chunk {
	var chunks []Chunk
	for i := 0; i < len(lines); i += c.opts.Size {
		end := i + c.opts.Size
		if end > len(lines) {
			end = len(lines)
		}
		chunks = append(chunks, Chunk{
			Label: fmt.Sprintf("%s-%d", c.opts.LabelPrefix, len(chunks)+1),
			Lines: append([]string(nil), lines[i:end]...),
		})
	}
	return chunks
}

func (c *Chunker) splitByDelimiter(lines []string) []Chunk {
	var chunks []Chunk
	var current []string
	for _, l := range lines {
		if c.delim.MatchString(l) && len(current) > 0 {
			chunks = append(chunks, Chunk{
				Label: fmt.Sprintf("%s-%d", c.opts.LabelPrefix, len(chunks)+1),
				Lines: current,
			})
			current = nil
		}
		current = append(current, l)
	}
	if len(current) > 0 {
		chunks = append(chunks, Chunk{
			Label: fmt.Sprintf("%s-%d", c.opts.LabelPrefix, len(chunks)+1),
			Lines: current,
		})
	}
	return chunks
}
