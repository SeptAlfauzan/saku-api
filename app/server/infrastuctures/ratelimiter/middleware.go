package ratelimiter

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/septalfauzan/saku-api/app/server/domain"
)

// NewMiddleware returns Fiber middleware enforcing rl per client IP.
func NewMiddleware(rl RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !rl.Allow(c.Context(), c.IP()) {
			t := strconv.Itoa(rl.RetryAfter())
			remaining := strconv.Itoa(rl.CheckRemainingTimeS(c.IP()))
			c.Set(fiber.HeaderRetryAfter, t)
			return c.Status(fiber.StatusTooManyRequests).JSON(domain.Error{
				Text: "rate limit exceeded, please retry again in " + remaining,
			})
		}
		return c.Next()
	}
}
