// Package ratelimit implements a per-key token-bucket rate limiter.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a token-bucket rate limit for a single key: it holds up
// to `burst` tokens, refilling at `rate` tokens per second, and Allow
// consumes one token per call if one is available.
type Limiter struct {
	mu         sync.Mutex
	rate       float64 // tokens per second
	burst      float64 // bucket capacity
	tokens     float64
	lastRefill time.Time
	now        func() time.Time // injectable clock, for tests
}

// NewLimiter constructs a Limiter starting with a full bucket.
func NewLimiter(rate, burst float64) *Limiter {
	return &Limiter{
		rate:       rate,
		burst:      burst,
		tokens:     burst,
		lastRefill: time.Now(),
		now:        time.Now,
	}
}

// Allow reports whether a request may proceed right now, consuming one
// token if so.
//
// TODO(candidate): implement. Refill tokens based on elapsed time since
// lastRefill (at `rate` tokens/sec, capped at `burst`), advance
// lastRefill, then consume one token if at least one is available.
func (l *Limiter) Allow() bool {
	panic("not implemented")
}

// MultiLimiter tracks one Limiter per key (e.g. per API token or requester
// ID), so each key is rate-limited independently of the others.
type MultiLimiter struct {
	mu       sync.Mutex
	rate     float64
	burst    float64
	limiters map[string]*Limiter
	lastSeen map[string]time.Time
	now      func() time.Time
}

// NewMultiLimiter constructs a MultiLimiter; every key gets its own bucket
// with the same rate/burst.
func NewMultiLimiter(rate, burst float64) *MultiLimiter {
	return &MultiLimiter{
		rate:     rate,
		burst:    burst,
		limiters: make(map[string]*Limiter),
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow reports whether the given key may proceed right now, creating a
// fresh bucket for keys seen for the first time.
//
// BUG(candidate): this compiles, and passes under light load, but has a
// real correctness bug — find and fix it. `go test -race` will NOT catch
// it; it's not a data race, every map access is already lock-protected.
// Think about what two goroutines calling Allow with the same brand-new
// key, at the same time, actually do to each other.
func (m *MultiLimiter) Allow(key string) bool {
	m.mu.Lock()
	lim, ok := m.limiters[key]
	m.mu.Unlock()

	if !ok {
		lim = NewLimiter(m.rate, m.burst)
		lim.now = m.now
		m.mu.Lock()
		m.limiters[key] = lim
		m.mu.Unlock()
	}

	m.mu.Lock()
	m.lastSeen[key] = m.now()
	m.mu.Unlock()

	return lim.Allow()
}

// EvictIdle removes every key whose last Allow call was more than idleFor
// ago. Returns the number of keys evicted.
//
// TODO(candidate): implement.
func (m *MultiLimiter) EvictIdle(idleFor time.Duration) int {
	panic("not implemented")
}
