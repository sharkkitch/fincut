// Package reader provides utilities for reading structured log files
// with support for byte-range offsets and line-based iteration.
//
// The LineReader type allows consumers to read lines from any io.ReadSeeker,
// optionally constraining reads to a specific byte range within the source.
// This is useful for efficiently processing large log files by seeking
// directly to relevant sections rather than scanning from the beginning.
//
// Example usage:
//
//	reader, err := reader.NewLineReader(0, 1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//	lines, err := reader.ReadFrom(f)
package reader
