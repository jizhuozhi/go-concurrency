package ratelimiter

import (
	"runtime"
	"sync/atomic"

	"github.com/jizhuozhi/go-concurrency/internal/utils"
)

type ShardingTokenBucketRateLimiter struct {
	index  uint64
	shards []*TokenBucketRateLimiter
}

func NewShardingTokenBucketRateLimiter(rateInSeconds int64, capacity int64) *ShardingTokenBucketRateLimiter {
	parallelism := int64(runtime.GOMAXPROCS(-1))
	parallelism = utils.MinInt64(parallelism, rateInSeconds)
	parallelism = utils.MinInt64(parallelism, capacity)

	shardingRate := rateInSeconds / parallelism
	shardingCapacity := capacity / parallelism
	remindRate := rateInSeconds % parallelism
	remindCapacity := capacity % parallelism

	rates := make([]int64, parallelism)
	capacities := make([]int64, parallelism)

	for i := int64(0); i < parallelism; i++ {
		rates[i] = shardingRate
		capacities[i] = shardingCapacity

		if i < remindRate {
			rates[i]++
		}
		if i < remindCapacity {
			capacities[i]++
		}
	}

	shards := make([]*TokenBucketRateLimiter, parallelism)
	for i := int64(0); i < parallelism; i++ {
		shards[i] = NewTokenBucketRateLimiter(rates[i], capacities[i])
	}
	return &ShardingTokenBucketRateLimiter{
		index:  0,
		shards: shards,
	}
}

func (r *ShardingTokenBucketRateLimiter) Allow() bool {
	index := atomic.AddUint64(&r.index, 1)
	return r.shards[index%uint64(len(r.shards))].Allow()
}
