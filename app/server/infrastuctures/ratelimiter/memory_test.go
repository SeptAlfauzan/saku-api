package ratelimiter

import (
	"context"
	"testing"
	"time"
)

func TestNewMemoryLimiterZeroWindowNoPanic(t *testing.T) {
	m := NewMemoryLimiter(5, 0)
	defer m.Stop()
	if !m.Allow(context.Background(), "1.2.3.4") {
		t.Fatal("expect allow on first request")
	}
}

func TestNewMemoryLimiterNegativeWindowNoPanic(t *testing.T) {
	m := NewMemoryLimiter(5, -time.Second)
	defer m.Stop()
	if !m.Allow(context.Background(), "1.2.3.4") {
		t.Fatal("expect allow on first request")
	}
}

func TestAllowWithinLimit(t *testing.T) {
	m := NewMemoryLimiter(2, time.Minute)
	defer m.Stop()
	ctx := context.Background()
	if !m.Allow(ctx, "k1") {
		t.Fatal("req1 should be allowed")
	}
	if !m.Allow(ctx, "k1") {
		t.Fatal("req2 should be allowed")
	}
	if m.Allow(ctx, "k1") {
		t.Fatal("req3 should be rejected (over limit)")
	}
}

func TestAllowKeyIsolation(t *testing.T) {
	m := NewMemoryLimiter(1, time.Minute)
	defer m.Stop()
	ctx := context.Background()
	if !m.Allow(ctx, "ip-a") {
		t.Fatal("ip-a first should be allowed")
	}
	if !m.Allow(ctx, "ip-b") {
		t.Fatal("ip-b first should be allowed (separate key)")
	}
}

func TestAllowWindowReset(t *testing.T) {
	m := NewMemoryLimiter(1, 30*time.Millisecond)
	defer m.Stop()
	ctx := context.Background()
	if !m.Allow(ctx, "k1") {
		t.Fatal("first should be allowed")
	}
	if m.Allow(ctx, "k1") {
		t.Fatal("second within window should be rejected")
	}
	time.Sleep(40 * time.Millisecond)
	if !m.Allow(ctx, "k1") {
		t.Fatal("after window elapses new request should be allowed")
	}
}

func TestCheckRemainingTimeS(t *testing.T) {
	m := NewMemoryLimiter(5, time.Minute)
	defer m.Stop()
	ctx := context.Background()
	m.Allow(ctx, "k1")
	rem := m.CheckRemainingTimeS("k1")
	if rem < 55 || rem > 60 {
		t.Fatalf("expected ~60s remaining, got %d", rem)
	}
	if m.CheckRemainingTimeS("unknown") != 0 {
		t.Fatal("unknown key should report 0 remaining")
	}
}