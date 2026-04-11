// Package segment implements log-line segmentation for fincut.
//
// A Segmenter partitions an ordered slice of log lines into labelled
// Segment values using one of two strategies:
//
//   - Count-based: lines are grouped into fixed-size windows identified
//     by a sequential label ("segment-1", "segment-2", …).
//
//   - Time-based: lines are grouped into windows of a given Duration by
//     parsing a timestamp extracted from each line. The segment label is
//     the formatted timestamp of the window's first record.
//
// Segments are consumed downstream by the diff and output packages to
// compare or render named slices of log output.
package segment
