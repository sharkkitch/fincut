// Package ratelimit implements a sliding-window rate limiter for log lines.
//
// It can be used to cap the throughput of a log processing pipeline,
// dropping lines that exceed a configured maximum count within a given
// time window. This is useful for suppressing log storms while still
// passing through a representative sample of high-volume events.
//
// Basic usage:
//
//	limiter, err := ratelimit.New(ratelimit.Options{
//		MaxLines: 100,
//		Window:   time.Second,
//	})
//	filtered := limiter.Apply(lines)
package ratelimit
