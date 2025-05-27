package otel

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes OpenTelemetry tracing with the given configuration
func InitTracer(
	ctx context.Context,
	serviceName string,
	traceEndpoint string,
	traceURLPath string,
	traceDeployEnv string,
	sampleRate float64,
) (trace.Tracer, error) {
	tracer := otel.Tracer(serviceName)

	// Configure trace exporter options
	traceOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(traceEndpoint),
		otlptracehttp.WithURLPath(traceURLPath),
	}

	// Create trace exporter
	traceExporter, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		logMsg := fmt.Sprintf(
			"starting OpenTelemetry trace exporter - Service_name: %s, trace_endpoint: %s, trace_url_path: %s",
			serviceName, traceEndpoint, traceURLPath,
		)
		return nil, errors.Wrapf(err, logMsg)
	}

	// Set up trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(sampleRate),
			),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				semconv.DeploymentEnvironmentKey.String(traceDeployEnv),
			),
		),
		sdktrace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracer, nil
}
