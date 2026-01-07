package middleware

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	maxTokens  int
	lastRefill time.Time
	mu         sync.Mutex
}

type RateLimiter struct {
	buckets sync.Map
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

type RateLimitResult struct {
	Allowed         bool
	RemainingTokens int
	RetryAfter      int // seconds
}

func (rl *RateLimiter) CheckLimit(userID string, maxTokens int, refillSeconds int) RateLimitResult {
	bucket := rl.getOrCreateBucket(userID, maxTokens)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	refillPeriod := time.Duration(refillSeconds) * time.Second

	// Refill tokens
	if now.Sub(bucket.lastRefill) >= refillPeriod {
		bucket.tokens = bucket.maxTokens
		bucket.lastRefill = now
	}

	// Check tokens
	if bucket.tokens > 0 {
		bucket.tokens--
		return RateLimitResult{
			Allowed:         true,
			RemainingTokens: bucket.tokens,
		}
	}

	retryAfter := int(refillPeriod.Seconds() - now.Sub(bucket.lastRefill).Seconds())

	return RateLimitResult{
		Allowed:         false,
		RemainingTokens: 0,
		RetryAfter:      retryAfter,
	}
}

func (rl *RateLimiter) getOrCreateBucket(userID string, maxTokens int) *TokenBucket {
	val, _ := rl.buckets.LoadOrStore(userID, &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		lastRefill: time.Now(),
	})
	return val.(*TokenBucket)
}

func (rl *RateLimiter) RemoveUser(userID string) {
	rl.buckets.Delete(userID)
}
