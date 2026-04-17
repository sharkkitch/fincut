// Package linefilter provides simple include/exclude line filtering using
// regular expressions. It is intended as a lightweight complement to the
// full regex pipeline in internal/filter, useful when callers need a
// standalone, stateless filter without pipeline overhead.
//
// Lines are first checked against exclude patterns — any match causes the
// line to be dropped immediately. Remaining lines are then checked against
// include patterns; if any include patterns are defined, a line must match
// at least one to pass through.
//
// If only exclude patterns are provided, all lines that do not match an
// exclude pattern are kept.
package linefilter
