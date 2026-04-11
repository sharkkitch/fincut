package cache

import "errors"

// Sentinel errors returned by cache operations.
var (
	// ErrInvalidCapacity is returned when cap < 1.
	ErrInvalidCapacity = errors.New("cache: capacity must be at least 1")

	// ErrInvalidTTL is returned when ttl <= 0.
	ErrInvalidTTL = errors.New("cache: TTL must be greater than zero")
)
