package ratelimiter

import (
	"sync"
	"time"

	"github.com/jizhuozhi/go-concurrency/internal/utils"
)

type TokenBucketRateLimiter struct {
	rateInSeconds int64
	capacity      int64

	mu                  sync.Mutex
	lastUpdateInSeconds int64
	tokens              int64
}

func NewTokenBucketRateLimiter(rateInSeconds int64, capacity int64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		rateInSeconds: rateInSeconds,
		capacity:      capacity,
		tokens:        capacity,
	}
}

func (r *TokenBucketRateLimiter) Allow() bool {
	now := time.Now().Unix()

	r.mu.Lock()
	defer r.mu.Unlock()

	escaped := utils.MaxInt64(now-r.lastUpdateInSeconds, 0)
	tokens := utils.MinInt64(r.tokens+escaped*r.rateInSeconds, r.capacity)

	if tokens <= 0 {
		return false
	}

	r.lastUpdateInSeconds = now
	r.tokens = tokens - 1

	return true
}
