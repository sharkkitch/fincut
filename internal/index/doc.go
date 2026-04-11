// Package index implements a line-offset indexer for fincut log files.
//
// It builds an in-memory index mapping each line number to its byte offset
// within the original source, enabling efficient random-access reads and
// precise byte-range extraction without re-scanning the entire file.
//
// Typical usage:
//
//	idx, err := index.Build(reader)
//	entry, err := idx.Lookup(lineNumber)
//	entries, err := idx.Range(start, end)
package index
