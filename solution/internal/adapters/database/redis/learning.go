package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"math"
)

type learningStorage struct {
	db *redis.Client
}

func NewLearningStorage(db *redis.Client) *learningStorage {
	return &learningStorage{
		db: db,
	}
}

func (s *learningStorage) SetR0(ctx context.Context, r0 float64) error {
	tracer := otel.Tracer("learning-storage")
	ctx, span := tracer.Start(ctx, "SetR0")
	defer span.End()

	return s.db.Set(ctx, "r0", r0, 0).Err()
}

func (s *learningStorage) GetR0(ctx context.Context) float64 {
	tracer := otel.Tracer("learning-storage")
	ctx, span := tracer.Start(ctx, "GetR0")
	defer span.End()

	r0, err := s.db.Get(ctx, "r0").Float64()
	if err != nil || math.IsNaN(r0) {
		_ = s.SetR0(ctx, 0.5)
		return 0.5
	}

	return r0
}
