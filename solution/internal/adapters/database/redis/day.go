package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"solution/internal/adapters/logger"
)

type dayStorage struct {
	db *redis.Client
}

func NewDayStorage(db *redis.Client) *dayStorage {
	return &dayStorage{
		db: db,
	}
}

func (s *dayStorage) SetDay(ctx context.Context, day int) error {
	tracer := otel.Tracer("day-storage")
	ctx, span := tracer.Start(ctx, "SetDay")
	defer span.End()

	logger.Log.Debugf("Setting day: %d", day)

	return s.db.Set(ctx, "day", day, 0).Err()
}

func (s *dayStorage) GetDay(ctx context.Context) (int, error) {
	tracer := otel.Tracer("day-storage")
	ctx, span := tracer.Start(ctx, "GetDay")
	defer span.End()

	return s.db.Get(ctx, "day").Int()
}
