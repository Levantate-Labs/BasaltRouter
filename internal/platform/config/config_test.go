package config_test

import (
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
)

func TestLoadValidConfig(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "local")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "./secrets/local.key")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Database.DSN == "" {
		t.Fatal("expected database DSN to be set")
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Fatalf("redis addr = %q, want localhost:6379", cfg.Redis.Addr)
	}
}

func TestValidateMissingDSN(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "local")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "./secrets/local.key")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing DSN")
	}
}

func TestValidateMissingRedisAddr(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "")
	t.Setenv("BASALT_ENCRYPTION_MODE", "local")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "./secrets/local.key")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing redis addr")
	}
}

func TestValidateEncryptionModeLocal(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "local")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing local key path")
	}
}

func TestValidateEncryptionModeKMS(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "kms")
	t.Setenv("BASALT_ENCRYPTION_KMS_KEY_ID", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing kms key id")
	}
}

func TestValidateTLSRequiresCertPaths(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "local")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "./secrets/local.key")
	t.Setenv("BASALT_SERVER_TLS_ENABLED", "true")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for missing TLS cert paths")
	}
}

func TestValidateInvalidEncryptionMode(t *testing.T) {
	t.Setenv("BASALT_DATABASE_DSN", "postgres://basalt:basalt@localhost:5432/basalt?sslmode=disable")
	t.Setenv("BASALT_REDIS_ADDR", "localhost:6379")
	t.Setenv("BASALT_ENCRYPTION_MODE", "invalid")
	t.Setenv("BASALT_ENCRYPTION_LOCAL_KEY_PATH", "./secrets/local.key")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid encryption mode")
	}
}
