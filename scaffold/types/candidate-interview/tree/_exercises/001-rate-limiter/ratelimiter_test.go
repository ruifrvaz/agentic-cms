package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock lets tests control time deterministically instead of sleeping.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func newFakeClock(start time.Time) *fakeClock {
	return &fakeClock{t: start}
}

func (c *fakeClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func TestLimiter_AllowsUpToBurstThenDenies(t *testing.T) {
	clock := newFakeClock(time.Now())
	l := NewLimiter(1, 3) // 1 token/sec, burst of 3
	l.now = clock.now
	l.lastRefill = clock.now()

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() to succeed on call %d (within burst)", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() to fail once the burst is exhausted")
	}
}

func TestLimiter_RefillsOverTime(t *testing.T) {
	clock := newFakeClock(time.Now())
	l := NewLimiter(2, 2) // 2 tokens/sec, burst of 2
	l.now = clock.now
	l.lastRefill = clock.now()

	if !l.Allow() || !l.Allow() {
		t.Fatal("expected both initial tokens to be available")
	}
	if l.Allow() {
		t.Fatal("expected the bucket to be empty")
	}

	clock.advance(500 * time.Millisecond) // 0.5s * 2/s = 1 token
	if !l.Allow() {
		t.Fatal("expected exactly one token to have refilled after 500ms")
	}
	if l.Allow() {
		t.Fatal("expected only one token to have refilled, not two")
	}
}

func TestLimiter_RefillCapsAtBurst(t *testing.T) {
	clock := newFakeClock(time.Now())
	l := NewLimiter(10, 2) // fast refill, small burst
	l.now = clock.now
	l.lastRefill = clock.now()

	l.Allow()
	l.Allow()
	clock.advance(10 * time.Second) // plenty of time to overfill if uncapped

	allowed := 0
	for i := 0; i < 5; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 2 {
		t.Fatalf("expected refill to cap at burst=2, got %d allowed calls", allowed)
	}
}

func TestMultiLimiter_KeysAreIsolated(t *testing.T) {
	m := NewMultiLimiter(1, 1)
	if !m.Allow("alice") {
		t.Fatal("alice's first call should be allowed")
	}
	if m.Allow("alice") {
		t.Fatal("alice's second immediate call should be denied")
	}
	if !m.Allow("bob") {
		t.Fatal("bob should have his own independent bucket")
	}
}

func TestMultiLimiter_ConcurrentNewKeyNeverExceedsBurst(t *testing.T) {
	const burst = 5
	const goroutines = 300
	const trials = 20 // repeat trials: a single lost-update race is timing-dependent, not guaranteed every run

	for trial := 0; trial < trials; trial++ {
		m := NewMultiLimiter(0, burst) // rate=0: no refill, isolates the bug from refill timing

		var allowed int64
		var wg sync.WaitGroup
		ready := make(chan struct{})
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				<-ready // start all goroutines at (as close to) the same instant as possible
				if m.Allow("shared-new-key") {
					atomic.AddInt64(&allowed, 1)
				}
			}()
		}
		close(ready)
		wg.Wait()

		if allowed > burst {
			t.Fatalf("trial %d: burst=%d but %d concurrent calls were allowed on the same brand-new key", trial, burst, allowed)
		}
	}
}

func TestMultiLimiter_EvictIdleRemovesOnlyStaleKeys(t *testing.T) {
	clock := newFakeClock(time.Now())
	m := NewMultiLimiter(1, 1)
	m.now = clock.now

	m.Allow("stale")
	clock.advance(time.Hour)
	m.Allow("fresh")

	evicted := m.EvictIdle(30 * time.Minute)
	if evicted != 1 {
		t.Fatalf("expected exactly 1 key evicted, got %d", evicted)
	}
	if _, stillThere := m.limiters["stale"]; stillThere {
		t.Fatal("stale key should have been evicted")
	}
	if _, stillThere := m.limiters["fresh"]; !stillThere {
		t.Fatal("fresh key should NOT have been evicted")
	}
}
