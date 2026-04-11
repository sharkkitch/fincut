package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MultiFileReader reads lines sequentially from one or more log files.
type MultiFileReader struct {
	paths []string
}

// NewMultiFileReader creates a MultiFileReader for the given file paths.
// It validates that paths is non-empty and that each file exists and is
// accessible before returning.
func NewMultiFileReader(paths []string) (*MultiFileReader, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("multifile reader: at least one file path is required")
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			return nil, fmt.Errorf("multifile reader: cannot access %q: %w", p, err)
		}
	}
	return &MultiFileReader{paths: paths}, nil
}

// Lines reads all lines from every file in order and returns them as a
// flat slice. Blank trailing newlines are handled gracefully.
func (m *MultiFileReader) Lines() ([]string, error) {
	var all []string
	for _, p := range m.paths {
		lines, err := readFileLines(p)
		if err != nil {
			return nil, fmt.Errorf("multifile reader: reading %q: %w", p, err)
		}
		all = append(all, lines...)
	}
	return all, nil
}

// readFileLines opens a single file and returns its non-empty lines.
func readFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Preserve all lines including empty ones that are not just trailing.
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Trim a single trailing empty string produced by a terminal newline.
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}
