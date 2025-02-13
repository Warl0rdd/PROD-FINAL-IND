package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
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
	logger.Log.Debugf("Setting day: %d", day)

	return s.db.Set(ctx, "day", day, 0).Err()
}

func (s *dayStorage) GetDay(ctx context.Context) (int, error) {
	return s.db.Get(ctx, "day").Int()
}
