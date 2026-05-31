package cache_test

import (
	"context"
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/cache"
	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
)

func TestNewRedisClientTLSEnabled(t *testing.T) {
	_, err := cache.NewRedisClient(config.RedisConfig{
		Addr:       "localhost:6379",
		TLSEnabled: true,
		TLSCAPath:  "/tmp/ca.pem",
	})
	if err == nil {
		t.Fatal("expected error when tls is enabled in phase 0")
	}
}

func TestNewRedisClientSuccess(t *testing.T) {
	client, err := cache.NewRedisClient(config.RedisConfig{
		Addr:     "localhost:6379",
		PoolSize: 5,
	})
	if err != nil {
		t.Fatalf("NewRedisClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("expected client")
	}
}

func TestPingNilClient(t *testing.T) {
	err := cache.Ping(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}
