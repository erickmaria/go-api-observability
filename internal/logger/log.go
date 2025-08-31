package logger

import (
	"context"
	"github/erickmaria/go-api-observability/internal/config"
	trace "github/erickmaria/go-api-observability/internal/traces"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc/credentials"
)

var Log *slog.Logger

func newStdOutLogExporter() (*stdoutlog.Exporter, error) {
	// return stdoutlog.New(stdoutlog.WithPrettyPrint())
	return stdoutlog.New()
}

func newOtelGrpcLogExporter(ctx context.Context) (*otlploggrpc.Exporter, error) {

	var insecure = config.GetBool("observability.otel.insecure")
	var secureOption otlploggrpc.Option
	if insecure {
		secureOption = otlploggrpc.WithInsecure()
	} else {
		secureOption = otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	var collectorEndpoint = config.GetSring("observability.otel.endpoint")
	return otlploggrpc.New(ctx, secureOption, otlploggrpc.WithEndpoint(collectorEndpoint))
}

// func newOtelHttpLogExporter(ctx context.Context) (*otlploghttp.Exporter, error) {
// 	var collectorEndpoint = config.GetSring("observability.otel.endpoint")
// 	return otlploghttp.New(ctx, otlploghttp.WithEndpoint(collectorEndpoint))
// }

func NewLogger(ctx context.Context) func() {

	var exporter log.Exporter
	var err error

	directToCollector := config.GetBool("OBSERVABILITY_LOGS_COLLECTOR")

	if directToCollector {
		exporter, err = newOtelGrpcLogExporter(ctx)
	} else {
		exporter, err = newStdOutLogExporter()
	}

	if err != nil {
		slog.Error("log exporter load with error:", err)
	}

	resources, err := trace.NewResources()
	if err != nil {
		slog.Error("trace resource load with error:", err)
	}

	// Create the logger provider
	provider := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(exporter),
		),
		log.WithResource(resources),
	)

	// setup log Bridge: https://opentelemetry.io/docs/languages/go/instrumentation/#log-bridge
	// Set the logger provider globally
	// global.SetLoggerProvider(provider)

	// Instantiate a new slog logger
	Log := otelslog.NewLogger(
		config.GetSring("observability.otel.service-name"),
		otelslog.WithLoggerProvider(provider),
	)
	slog.SetDefault(Log)
	slog.Info("logger setup successfully")

	return func() {
		err := provider.Shutdown(ctx)
		slog.Error("log shutdown:", err)
	}
}
