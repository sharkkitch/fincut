package reader

import (
	"fmt"
	"io"
	"os"
)

// MultiFileReader reads lines sequentially from multiple files,
// tagging each line with its source filename.
type MultiFileReader struct {
	paths []string
}

// TaggedLine holds a line of text along with the file it originated from.
type TaggedLine struct {
	File string
	Line string
}

// NewMultiFileReader creates a MultiFileReader for the given file paths.
// Returns an error if no paths are provided.
func NewMultiFileReader(paths []string) (*MultiFileReader, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("multifilereader: at least one file path is required")
	}
	return &MultiFileReader{paths: paths}, nil
}

// ReadAll reads all lines from all files in order, returning tagged lines.
// If a file cannot be opened or read, an error is returned immediately.
func (m *MultiFileReader) ReadAll() ([]TaggedLine, error) {
	var results []TaggedLine

	for _, path := range m.paths {
		lines, err := readFileLines(path)
		if err != nil {
			return nil, fmt.Errorf("multifilereader: reading %q: %w", path, err)
		}
		for _, line := range lines {
			results = append(results, TaggedLine{File: path, Line: line})
		}
	}

	return results, nil
}

// readFileLines opens a file and returns all non-empty lines.
func readFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			line := string(data[start:i])
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}
	return lines, nil
}
