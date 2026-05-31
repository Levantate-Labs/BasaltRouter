package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"github.com/LevantateLabs/basaltrouter/internal/platform/log"
)

func TestInitJSONLogger(t *testing.T) {
	logger := log.Init(config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	})
	if logger == nil {
		t.Fatal("expected logger")
	}
}

func TestFromContextWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	ctx := log.WithContext(context.Background(), logger)
	ctx = log.WithRequestID(ctx, "req-123")

	log.FromContext(ctx).Info("test event")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log: %v", err)
	}
	if entry["request_id"] != "req-123" {
		t.Fatalf("request_id = %v, want req-123", entry["request_id"])
	}
}

func TestFromContextFallback(t *testing.T) {
	logger := log.FromContext(context.Background())
	if logger == nil {
		t.Fatal("expected default logger")
	}
}

func TestInitTextLogger(t *testing.T) {
	logger := log.Init(config.LoggingConfig{
		Level:  "debug",
		Format: "text",
		Output: "stderr",
	})
	if logger == nil {
		t.Fatal("expected logger")
	}
}
