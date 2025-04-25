package ratelimiter

import (
	"math"
	"testing"
)

func BenchmarkTokenBucketRateLimiter_Allow(b *testing.B) {
	limiter := NewTokenBucketRateLimiter(math.MaxInt64, math.MaxInt64)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}

func BenchmarkShardingTokenBucketRateLimiter_Allow(b *testing.B) {
	limiter := NewShardingTokenBucketRateLimiter(math.MaxInt64, math.MaxInt64)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}
