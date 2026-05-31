package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/database"
	"github.com/LevantateLabs/basaltrouter/internal/platform/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: load config: %v\n", err)
		os.Exit(1)
	}

	logger := log.Init(cfg.Logging)

	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Error("seed: connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	const insertSQL = `
		INSERT INTO audit.audit_log (action, resource_type, resource_id, metadata)
		VALUES ($1, $2, $3, $4)
	`

	_, err = pool.Exec(
		ctx, insertSQL,
		"seed.connectivity_check",
		"system",
		"seed",
		`{"source":"cmd/seed"}`,
	)
	if err != nil {
		logger.Error("seed: insert audit row", "error", err)
		os.Exit(1)
	}

	logger.Info("seed completed successfully")
}
