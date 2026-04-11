// Package cache provides a simple in-memory line cache for recently read
// log segments, reducing repeated disk I/O during interactive filtering.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached slice of lines along with metadata.
type Entry struct {
	Lines     []string
	CachedAt  time.Time
	ByteStart int64
	ByteEnd   int64
}

// Cache is a thread-safe, capacity-bounded store keyed by file path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	cap     int
	ttl     time.Duration
}

// New creates a Cache with the given capacity and TTL.
// Returns an error if cap < 1 or ttl <= 0.
func New(cap int, ttl time.Duration) (*Cache, error) {
	if cap < 1 {
		return nil, ErrInvalidCapacity
	}
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Cache{
		entries: make(map[string]*Entry, cap),
		cap:     cap,
		ttl:     ttl,
	}, nil
}

// Set stores lines under the given key, evicting the oldest entry when at
// capacity.
func (c *Cache) Set(key string, e *Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.cap {
		c.evictOldest()
	}
	e.CachedAt = time.Now()
	c.entries[key] = e
}

// Get retrieves an entry. Returns nil, false when missing or expired.
func (c *Cache) Get(key string) (*Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Since(e.CachedAt) > c.ttl {
		return nil, false
	}
	return e, true
}

// Delete removes a single entry.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the current number of entries (including possibly expired ones).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// evictOldest removes the entry with the oldest CachedAt timestamp.
// Caller must hold the write lock.
func (c *Cache) evictOldest() {
	var oldest string
	var oldestTime time.Time
	for k, e := range c.entries {
		if oldest == "" || e.CachedAt.Before(oldestTime) {
			oldest = k
			oldestTime = e.CachedAt
		}
	}
	if oldest != "" {
		delete(c.entries, oldest)
	}
}
