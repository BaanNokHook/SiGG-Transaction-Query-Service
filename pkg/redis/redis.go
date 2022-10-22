package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisClient(opts ...Option) *RedisCache {
	rdb := &redis.Options{}

	//Custom Options
	for _, opts := range opts {
		opts(rdb)
	}

	return &RedisCache{
		redisClient: redis.NewClient(rdb),
	}

}

func (r *RedisCache) Get(key string) (string, error) {
	var ctx = context.Background()
	val, err := r.redisClient.Get(ctx, key).Result()
	return val, err
}

// Close -.
func (r *RedisCache) Close() error {
	//Close connection here
	return r.redisClient.Close()
}
