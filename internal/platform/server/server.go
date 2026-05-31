package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LevantateLabs/basaltrouter/internal/platform/cache"
	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/database"
	"github.com/LevantateLabs/basaltrouter/internal/platform/health"
	"github.com/LevantateLabs/basaltrouter/internal/platform/log"
	"github.com/LevantateLabs/basaltrouter/internal/platform/middleware"
	"github.com/LevantateLabs/basaltrouter/internal/platform/telemetry"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Resources struct {
	Config      *config.Config
	Logger      *slog.Logger
	Pool        *pgxpool.Pool
	Redis       *redis.Client
	ShutdownTel func(context.Context) error
}

func Bootstrap(ctx context.Context, serviceName string) (*Resources, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	cfg.Telemetry.ServiceName = serviceName

	logger := log.Init(cfg.Logging)

	shutdownTel, err := telemetry.Init(ctx, cfg.Telemetry)
	if err != nil {
		return nil, err
	}

	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		_ = shutdownTel(ctx)
		return nil, err
	}

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		pool.Close()
		_ = shutdownTel(ctx)
		return nil, err
	}

	if err := cache.Ping(ctx, redisClient); err != nil {
		_ = redisClient.Close()
		pool.Close()
		_ = shutdownTel(ctx)
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return &Resources{
		Config:      cfg,
		Logger:      logger,
		Pool:        pool,
		Redis:       redisClient,
		ShutdownTel: shutdownTel,
	}, nil
}

func (r *Resources) Close(ctx context.Context) {
	if r.Redis != nil {
		_ = r.Redis.Close()
	}
	if r.Pool != nil {
		r.Pool.Close()
	}
	if r.ShutdownTel != nil {
		_ = r.ShutdownTel(ctx)
	}
}

func NewHTTPServer(res *Resources, serviceName string) *http.Server {
	mux := http.NewServeMux()

	checkers := []health.Checker{
		health.NewFuncChecker("postgres", func(ctx context.Context) error {
			return database.Ping(ctx, res.Pool)
		}),
		health.NewFuncChecker("redis", func(ctx context.Context) error {
			return cache.Ping(ctx, res.Redis)
		}),
	}
	health.NewHandler(checkers...).Register(mux)

	handler := middleware.RequestID(telemetry.HTTPMiddleware(serviceName)(mux))

	addr := fmt.Sprintf(":%d", res.Config.Server.Port)
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  res.Config.Server.ReadTimeout,
		WriteTimeout: res.Config.Server.WriteTimeout,
	}
}

func RunHTTPServer(ctx context.Context, res *Resources, serviceName string) error {
	srv := NewHTTPServer(res, serviceName)

	errCh := make(chan error, 1)
	go func() {
		res.Logger.Info("server starting", "addr", srv.Addr, "service", serviceName)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		res.Logger.Info("shutdown signal received", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	res.Logger.Info("server stopped", "service", serviceName)
	return nil
}
