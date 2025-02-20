package logger

import (
	"context"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log         *logger
	LogProvider *sdklog.LoggerProvider
)

type logger struct {
	*zap.SugaredLogger
}

// New is a function to initialize logger
/*
 * debug bool - is debug mode
 */
func New(debug bool) {
	var level zapcore.Level
	if debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	provider, err := newLoggerProvider(context.Background())

	if err != nil {
		panic(err)
	}

	core := otelzap.NewCore("internal/adapters/logger", otelzap.WithLoggerProvider(provider))

	core.Enabled(level)

	log := zap.New(core, zap.AddCaller())

	Log = &logger{
		SugaredLogger: log.Sugar(),
	}

	LogProvider = provider
}

// Да, плохо, но пришлось эту функцию перенести сюда из app.go

func newLoggerProvider(ctx context.Context) (*sdklog.LoggerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Fiber-Backend"),
		),
	)

	if err != nil {
		return nil, err
	}

	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithEndpoint("otel-collector:4318"), otlploghttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(r),
	)
	return loggerProvider, nil
}
