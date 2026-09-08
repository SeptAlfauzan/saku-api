package ratelimiter

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
	RetryAfter() int
	CheckRemainingTimeS(key string) int
}
