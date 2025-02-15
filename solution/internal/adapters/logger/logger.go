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
	"os"
	"time"
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
 * timeZone string - logger time zone, by default "GMT"
 */
func New(debug bool, timeZone string) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Цветная подсветка уровней
		EncodeTime:     customTimeEncoder,                // Кастомный формат времени
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Краткий формат caller
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	if timeZone != "" {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(time.FixedZone(timeZone, 3*60*60)).Format("2006-01-02 15:04:05"))
		}
	}

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

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
		otelzap.NewCore("internal/adapters/logger", otelzap.WithLoggerProvider(provider)),
	)

	log := zap.New(core, zap.AddCaller())

	Log = &logger{
		SugaredLogger: log.Sugar(),
	}

	LogProvider = provider
}

// customTimeEncoder форматирует время в GMT+0
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02 15:04:05"))
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
