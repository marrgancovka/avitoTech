package repo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func (cache *RedisRepo) GetBanner(ctx context.Context, key string) ([]byte, bool) {
	value, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return value, true
}
func (cache *RedisRepo) SetBanner(ctx context.Context, key string, data []byte) error {
	return cache.client.Set(ctx, key, data, time.Minute*5).Err()
}
