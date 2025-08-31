package trace

import (
	"context"
	"fmt"
	"github/erickmaria/go-api-observability/internal/config"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"
)

// https://betterstack.com/community/guides/observability/opentelemetry-go
var Tracer = otel.Tracer("server")

func newStdOutTraceExporter() (*stdouttrace.Exporter, error) {
	// return stdouttrace.New(stdouttrace.WithPrettyPrint())
	return stdouttrace.New()
}

func newOtelTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	var collectorEndpoint = config.GetSring("observability.otel.endpoint")
	var insecure = config.GetBool("observability.otel.insecure")
	var secureOption otlptracegrpc.Option

	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	} else {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	return otlptrace.New(ctx, otlptracegrpc.NewClient(
		secureOption,
		otlptracegrpc.WithEndpoint(collectorEndpoint),
	))

}

func InitTracer(ctx context.Context) func() {

	var exporter sdktrace.SpanExporter
	var err error

	// exporter, err := newStdOutTraceExporter()
	// exporter, err := newOtelTraceExporter(ctx)

	directToCollector := config.GetBool("observability.trace.collector")
	if directToCollector {
		exporter, err = newOtelTraceExporter(ctx)
	} else {
		exporter, err = newStdOutTraceExporter()
	}

	if err != nil {
		slog.Error("trace exporter load with error:", err)
	}

	resources, err := NewResources()
	if err != nil {
		slog.Error("trace resource load with error:", err)
	}

	provider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	slog.Info("trace setup successfully")

	return func() {
		err := provider.Shutdown(ctx)
		slog.Error("trace shutdown:", err)
	}
}

func Middleware(next http.Handler) http.Handler {
	return middleware(next)
}

func middleware(handler http.Handler, opts ...otelhttp.Option) http.Handler {

	httpSpanName := func(operation string, r *http.Request) string {
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
	}

	opts = append(opts, otelhttp.WithSpanNameFormatter(httpSpanName))

	return otelhttp.NewHandler(handler, "/", opts...)
}

func NewResources() (*resource.Resource, error) {
	var serviceName = config.GetSring("observability.service-name")
	resources, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		slog.Error("trace resource load with error:", err)
	}

	return resources, err
}
