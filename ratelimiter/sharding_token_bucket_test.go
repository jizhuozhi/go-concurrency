package ratelimiter

import (
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestShardingTokenBucketRateLimiter_Allow(t *testing.T) {
	PatchConvey("TestShardingTokenBucketRateLimiter_Allow", t, func() {
		origin := time.Now
		base := time.Now().Round(time.Second)
		Mock(time.Now).To(func() time.Time {
			return base
		}).Origin(&origin).Build()

		assert.Equal(t, base, time.Now())
		rate := int64(10)
		capacity := int64(20)
		ratelimiter := NewShardingTokenBucketRateLimiter(rate, capacity)

		// allow burst requests with capacity
		for i := int64(0); i < capacity; i++ {
			assert.True(t, ratelimiter.Allow())
		}
		for i := int64(0); i < capacity; i++ {
			assert.False(t, ratelimiter.Allow())
		}

		// allow requests with rate
		base = base.Add(1 * time.Second)
		for i := int64(0); i < rate; i++ {
			assert.True(t, ratelimiter.Allow())
		}
		for i := int64(0); i < capacity; i++ {
			assert.False(t, ratelimiter.Allow())
		}

		// allow requests with at most capacity
		base = base.Add(time.Duration(capacity/rate*2) * time.Second)
		for i := int64(0); i < capacity; i++ {
			assert.True(t, ratelimiter.Allow())
		}
		for i := int64(0); i < capacity; i++ {
			assert.False(t, ratelimiter.Allow())
		}
	})
}
