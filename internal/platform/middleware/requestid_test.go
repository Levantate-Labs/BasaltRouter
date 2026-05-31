package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LevantateLabs/basaltrouter/internal/platform/middleware"
)

func TestRequestIDGenerated(t *testing.T) {
	var captured string
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = middleware.RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if captured == "" {
		t.Fatal("expected request ID in context")
	}
	if rec.Header().Get(middleware.RequestIDHeader) != captured {
		t.Fatalf("header %q != context %q", rec.Header().Get(middleware.RequestIDHeader), captured)
	}
}

func TestRequestIDPropagatedFromHeader(t *testing.T) {
	const existing = "550e8400-e29b-41d4-a716-446655440000"

	var captured string
	handler := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = middleware.RequestIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middleware.RequestIDHeader, existing)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if captured != existing {
		t.Fatalf("context id = %q, want %q", captured, existing)
	}
	if rec.Header().Get(middleware.RequestIDHeader) != existing {
		t.Fatalf("header id = %q, want %q", rec.Header().Get(middleware.RequestIDHeader), existing)
	}
}

func TestRequestIDFromContextEmpty(t *testing.T) {
	if id := middleware.RequestIDFromContext(context.Background()); id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
}
