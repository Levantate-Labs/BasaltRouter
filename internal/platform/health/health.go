package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

type Handler struct {
	checkers []Checker
}

func NewHandler(checkers ...Checker) *Handler {
	return &Handler{checkers: checkers}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.Liveness)
	mux.HandleFunc("GET /readyz", h.Readiness)
}

func (h *Handler) Liveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	failures := h.runChecks(ctx)

	if len(failures) > 0 {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status":   "not_ready",
			"failures": failures,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h *Handler) runChecks(ctx context.Context) map[string]string {
	failures := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, checker := range h.checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()
			if err := c.Check(ctx); err != nil {
				mu.Lock()
				failures[c.Name()] = err.Error()
				mu.Unlock()
			}
		}(checker)
	}

	wg.Wait()
	return failures
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type FuncChecker struct {
	name string
	fn   func(ctx context.Context) error
}

func NewFuncChecker(name string, fn func(ctx context.Context) error) *FuncChecker {
	return &FuncChecker{name: name, fn: fn}
}

func (c *FuncChecker) Name() string {
	return c.name
}

func (c *FuncChecker) Check(ctx context.Context) error {
	return c.fn(ctx)
}
