// controller/main.go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/tree/main/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	// Initialize distributed tracing via Jaeger.
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	// Ensure tracer provider shuts down when the application exits.
	defer func() { _ = tp.Shutdown(context.Background()) }()

	// Register application endpoints.
	registerHandlers()

	// Expose Prometheus metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())

	log.Println("ChaosLab Controller running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Controller failed to start: %v", err)
	}
}

// initTracer sets up an OpenTelemetry tracer provider with Jaeger exporter.
func initTracer() (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://jaeger-collector:14268/api/traces"),
	))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
