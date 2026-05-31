package telemetry_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/telemetry"
)

func TestInitDisabled(t *testing.T) {
	shutdown, err := telemetry.Init(context.Background(), config.TelemetryConfig{
		Enabled:     false,
		ServiceName: "test",
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error = %v", err)
	}
}

func TestInitStdoutExporter(t *testing.T) {
	shutdown, err := telemetry.Init(context.Background(), config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "test",
		Exporter:    "stdout",
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	if telemetry.Tracer() == nil {
		t.Fatal("expected tracer")
	}
}

func TestInitUnsupportedExporter(t *testing.T) {
	_, err := telemetry.Init(context.Background(), config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "test",
		Exporter:    "otlp",
	})
	if err == nil {
		t.Fatal("expected error for unsupported exporter")
	}
}

func TestHTTPMiddleware(t *testing.T) {
	shutdown, err := telemetry.Init(context.Background(), config.TelemetryConfig{
		Enabled:     true,
		ServiceName: "test",
		Exporter:    "stdout",
	})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	called := false
	handler := telemetry.HTTPMiddleware("test")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected handler to be called")
	}
}
