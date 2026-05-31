package server_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/server"
)

func testResources() *server.Resources {
	return &server.Resources{
		Config: &config.Config{
			Server: config.ServerConfig{
				Port:         8080,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 5 * time.Second,
			},
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestNewHTTPServerHealthRoutes(t *testing.T) {
	res := testResources()
	srv := server.NewHTTPServer(res, "test")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want 200", rec.Code)
	}
}

func TestNewHTTPServerReadyzWithoutDeps(t *testing.T) {
	res := testResources()
	srv := server.NewHTTPServer(res, "test")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("readyz status = %d, want 503 without deps", rec.Code)
	}
}

func TestResourcesClose(t *testing.T) {
	res := testResources()
	res.Close(context.Background())
}
