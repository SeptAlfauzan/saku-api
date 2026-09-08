package ratelimiter

import (
	"context"
	"math"
	"sync"
	"time"
)

type Bucket struct {
	count       int
	windowStart time.Time
}

type MemoryLimiter struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
	limit   int
	window  time.Duration
	stopCh  chan struct{}
	close   sync.Once
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	m := &MemoryLimiter{
		buckets: make(map[string]*Bucket),
		limit:   limit,
		window:  window,
		stopCh:  make(chan struct{}),
	}
	go m.sweep()
	return m
}

func (m *MemoryLimiter) Allow(ctx context.Context, key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	b, exists := m.buckets[key]
	if !exists || now.Sub(b.windowStart) > m.window {
		m.buckets[key] = &Bucket{count: 1, windowStart: now}
		return true
	}
	b.count++
	return b.count <= m.limit
}

// RetryAfter returns the window length in whole seconds, used for the
// Retry-After header when a request is rejected.
func (m *MemoryLimiter) RetryAfter() int {
	return int(m.window.Seconds())
}

func (m *MemoryLimiter) CheckRemainingTimeS(key string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.buckets[key]
	if !ok {
		return 0
	}

	rem := m.window - time.Since(b.windowStart)

	if rem < 0 {
		rem = 0
	}

	return int(math.Ceil(rem.Seconds()))
}

func (m *MemoryLimiter) Stop() {
	m.close.Do(func() { close(m.stopCh) })
}

func (m *MemoryLimiter) sweep() {
	t := time.NewTicker(m.window)
	defer t.Stop()
	for {
		select {
		case <-m.stopCh:
			return
		case <-t.C:
			m.mu.Lock()
			now := time.Now()
			for k, b := range m.buckets {
				if now.Sub(b.windowStart) > m.window {
					delete(m.buckets, k)
				}
			}
			m.mu.Unlock()
		}
	}
}
