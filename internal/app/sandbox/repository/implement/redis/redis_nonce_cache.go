package repository_implement

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	repoInterface "github.com/vando2108/sandbox_service/internal/app/sandbox/repository/interface"
)

type RedisNonceCache struct {
	client *redis.Client
}

func NewRedisNonceCache(client *redis.Client) repoInterface.NonceCacheRepository {
	return &RedisNonceCache{
		client: client,
	}
}

func (r *RedisNonceCache) SetNonce(ctx context.Context, publickey string, nonce string, expireTime time.Duration) error {
	if err := r.client.Set(ctx, publickey, nonce, expireTime).Err(); err != nil {
		return fmt.Errorf("failed to set nonce: %w", err)
	}

	return nil
}

func (r *RedisNonceCache) GetNonce(ctx context.Context, publickey string) (string, error) {
	nonce, err := r.client.Get(ctx, publickey).Result()
	return nonce, err
}
