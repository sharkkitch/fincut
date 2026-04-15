// Package limit implements a hard-stop limiter for line and byte budgets.
//
// A Limiter is constructed with New and accepts Options that specify either a
// maximum number of lines, a maximum number of bytes, or both.  When both are
// provided the first ceiling reached wins.
//
// Typical use:
//
//	lim, err := limit.New(limit.Options{MaxLines: 1000})
//	if err != nil { ... }
//	result := lim.Apply(inputLines)
package limit
