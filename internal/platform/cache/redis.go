package cache

import (
	"context"
	"fmt"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
	}

	if cfg.TLSEnabled {
		return nil, fmt.Errorf("redis: tls is configured but not implemented in phase 0")
	}

	client := redis.NewClient(opts)
	return client, nil
}

func Ping(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return client.Ping(ctx).Err()
}
