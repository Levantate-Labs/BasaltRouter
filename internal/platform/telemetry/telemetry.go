package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func Init(ctx context.Context, cfg config.TelemetryConfig) (func(context.Context) error, error) {
	if !cfg.Enabled {
		tracer = otel.Tracer(cfg.ServiceName)
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("telemetry: create stdout exporter: %w", err)
	}

	if strings.ToLower(cfg.Exporter) != "stdout" {
		return nil, fmt.Errorf("telemetry: unsupported exporter %q (only stdout in phase 0)", cfg.Exporter)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("telemetry: create resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = provider.Tracer(cfg.ServiceName)

	return provider.Shutdown, nil
}

func Tracer() trace.Tracer {
	if tracer != nil {
		return tracer
	}
	return otel.Tracer("basaltrouter")
}

func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			spanName := r.Method + " " + r.URL.Path
			ctx, span := Tracer().Start(
				ctx, spanName,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.route", r.URL.Path),
					attribute.String("service.name", serviceName),
				),
			)
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
