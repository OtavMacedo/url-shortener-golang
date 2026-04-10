package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) IsRateLimited(ctx context.Context, ip string) (bool, error) {
	key := "ratelimit:" + ip
	now := time.Now()
	windowStart := now.Add(-1 * time.Minute).UnixMilli()

	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	count := pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: now.UnixMilli(),
	})
	pipe.Expire(ctx, key, time.Minute)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return count.Val() >= 10, nil
}
