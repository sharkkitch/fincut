// Package index provides line-offset indexing for fast byte-range lookups
// within structured log files processed by fincut.
package index

import (
	"bufio"
	"fmt"
	"io"
)

// Entry records the byte offset and line number of a single log line.
type Entry struct {
	Line   int
	Offset int64
	Length int
}

// Index maps line numbers to their byte offsets in a source stream.
type Index struct {
	entries []Entry
}

// Build reads from r and constructs a byte-offset index for every line.
func Build(r io.Reader) (*Index, error) {
	if r == nil {
		return nil, fmt.Errorf("index: reader must not be nil")
	}

	idx := &Index{}
	scanner := bufio.NewScanner(r)
	var offset int64
	lineNum := 0

	for scanner.Scan() {
		raw := scanner.Bytes()
		length := len(raw) + 1 // +1 for newline
		idx.entries = append(idx.entries, Entry{
			Line:   lineNum,
			Offset: offset,
			Length: length,
		})
		offset += int64(length)
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("index: scan error: %w", err)
	}

	return idx, nil
}

// Lookup returns the Entry for the given line number, or an error if out of range.
func (idx *Index) Lookup(line int) (Entry, error) {
	if line < 0 || line >= len(idx.entries) {
		return Entry{}, fmt.Errorf("index: line %d out of range [0, %d)", line, len(idx.entries))
	}
	return idx.entries[line], nil
}

// Len returns the total number of indexed lines.
func (idx *Index) Len() int {
	return len(idx.entries)
}

// Range returns all entries between startLine and endLine (inclusive).
func (idx *Index) Range(startLine, endLine int) ([]Entry, error) {
	if startLine < 0 || endLine >= len(idx.entries) || startLine > endLine {
		return nil, fmt.Errorf("index: invalid range [%d, %d] for index of length %d",
			startLine, endLine, len(idx.entries))
	}
	return idx.entries[startLine : endLine+1], nil
}
