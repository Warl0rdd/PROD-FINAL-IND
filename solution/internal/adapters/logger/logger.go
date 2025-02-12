package logger

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	Log *logger
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

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	log := zap.New(core, zap.AddCaller())

	Log = &logger{
		SugaredLogger: log.Sugar(),
	}
}

// customTimeEncoder форматирует время в GMT+0
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.In(time.FixedZone("GMT+0", 3*60*60)).Format("2006-01-02 15:04:05"))
}

// Адаптер для pgx

type ZapQueryTracer struct {
	Log *zap.SugaredLogger
}

// TraceQueryStart вызывается в начале выполнения запроса.
func (t *ZapQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	t.Log.Infow("Начало выполнения запроса",
		"sql", data.SQL,
		"args", data.Args,
	)
	return ctx
}

// TraceQueryEnd вызывается по окончании выполнения запроса.
func (t *ZapQueryTracer) TraceQueryEnd(
	_ context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	if data.Err != nil {
		t.Log.Errorw("Запрос завершился с ошибкой",
			"error", data.Err,
		)
	} else {
		t.Log.Infow("Запрос успешно выполнен")
	}
}
