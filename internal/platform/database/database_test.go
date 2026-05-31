package database_test

import (
	"context"
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/database"
)

func TestNewPoolInvalidDSN(t *testing.T) {
	_, err := database.NewPool(context.Background(), config.DatabaseConfig{
		DSN: "not-a-valid-dsn",
	})
	if err == nil {
		t.Fatal("expected error for invalid dsn")
	}
}

func TestPingNilPool(t *testing.T) {
	err := database.Ping(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil pool")
	}
}
