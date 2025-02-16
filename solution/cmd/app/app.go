package app

import (
	"context"
	"errors"
	"github.com/FKouhai/fiberprometheus/v3"
	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sheeiavellie/go-yandexgpt"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"os"
	"os/signal"
	"solution/internal/adapters/config"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/logger"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// App is a struct that contains the fiber app, database connection, listen port, validator, logging boolean etc.
type App struct {
	Fiber     *fiber.App
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Minio     *minio.Client
	Validator *validator.Validator
	GPT       *yandexgpt.YandexGPTClient
}

// New is a function that creates a new app struct
func New(config *config.Config) *App {
	fiberApp := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(validator.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	},
	)

	fiberApp.Use(recoverer.New(
		recoverer.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c fiber.Ctx, e any) {
				logger.Log.Error(e)
			},
		}))

	prom := fiberprometheus.New("Backend")
	prom.RegisterAt(fiberApp, "/metrics")
	fiberApp.Use(prom.Middleware)

	return &App{
		Fiber:     fiberApp,
		DB:        config.Database,
		Redis:     config.Redis,
		Minio:     config.Minio,
		Validator: validator.New(),
		GPT:       config.GPT,
	}
}

// Start is a function that starts the app
func (a *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		logger.Log.Panicf("failed to setup otel: %v", err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	if viper.GetBool("settings.listen-tls") {
		if err := a.Fiber.Listen(
			":"+os.Getenv("SERVER_PORT"),
			fiber.ListenConfig{
				CertFile:    viper.GetString("service.backend.certificate.cert-file"),
				CertKeyFile: viper.GetString("service.backend.certificate.key-file"),
			}); err != nil {
			logger.Log.Panicf("failed to start listen (with tls): %v", err)
		}
	} else {
		logger.Log.Debugf("port: %s", viper.GetString("service.backend.port"))
		if err := a.Fiber.Listen(":" + os.Getenv("SERVER_PORT")); err != nil {
			logger.Log.Panicf("failed to start listen (no tls): %v", err)
		}
	}
}

func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up logger provider.
	loggerProvider := logger.LogProvider

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("Fiber-Backend"),
		),
	)

	traceExporter, err := otlptracehttp.New(
		ctx, otlptracehttp.WithEndpoint("otel-collector:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second*5)),
		trace.WithResource(r),
	)
	return traceProvider, nil
}
