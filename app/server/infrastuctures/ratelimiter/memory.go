package ratelimiter

import (
	"context"
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
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{
		buckets: make(map[string]*Bucket),
		limit:   limit,
		window:  window,
	}
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
