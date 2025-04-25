package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter_Allow(t *testing.T) {
	rate := int64(10)
	capacity := int64(20)
	ratelimiter := NewTokenBucketRateLimiter(rate, capacity)

	// allow burst requests with capacity
	for i := int64(0); i < capacity; i++ {
		assert.True(t, ratelimiter.Allow())
	}
	for i := int64(0); i < capacity; i++ {
		assert.False(t, ratelimiter.Allow())
	}

	// allow requests with rate
	time.Sleep(1 * time.Second)
	for i := int64(0); i < rate; i++ {
		assert.True(t, ratelimiter.Allow())
	}
	for i := int64(0); i < capacity; i++ {
		assert.False(t, ratelimiter.Allow())
	}

	// allow requests with at most capacity
	time.Sleep(time.Duration(capacity/rate*2) * time.Second)
	for i := int64(0); i < capacity; i++ {
		assert.True(t, ratelimiter.Allow())
	}
	for i := int64(0); i < capacity; i++ {
		assert.False(t, ratelimiter.Allow())
	}
}
