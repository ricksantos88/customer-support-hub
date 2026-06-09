package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

type windowEntry struct {
	count     int
	windowEnd time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	windows map[string]*windowEntry
	limit   int
	window  time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		windows: make(map[string]*windowEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, ok := rl.windows[key]
	if !ok || now.After(entry.windowEnd) {
		rl.windows[key] = &windowEntry{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

// cleanup removes expired entries every minute to prevent memory leaks.
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.windows {
			if now.After(entry.windowEnd) {
				delete(rl.windows, key)
			}
		}
		rl.mu.Unlock()
	}
}

func tooManyRequests(c fiber.Ctx, retryAfter time.Duration) error {
	c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many requests"})
}

// NewIPRateLimiter limits requests per client IP address.
func NewIPRateLimiter(limitPerMinute int) fiber.Handler {
	rl := newRateLimiter(limitPerMinute, time.Minute)
	return func(c fiber.Ctx) error {
		if !rl.allow(c.IP()) {
			return tooManyRequests(c, time.Minute)
		}
		return c.Next()
	}
}

// NewAgentRateLimiter limits requests per authenticated agent (agent_id JWT claim).
// Must be applied after the BearerAuthMiddleware.
func NewAgentRateLimiter(limitPerMinute int) fiber.Handler {
	rl := newRateLimiter(limitPerMinute, time.Minute)
	return func(c fiber.Ctx) error {
		agentID, ok := c.Locals("agent_id").(string)
		if !ok || agentID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing agent context"})
		}
		if !rl.allow(agentID) {
			return tooManyRequests(c, time.Minute)
		}
		return c.Next()
	}
}
