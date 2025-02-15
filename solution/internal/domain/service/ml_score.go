package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
	"solution/internal/domain/utils/ads"
)

type mlScoreStorage interface {
	InsertOrUpdateMlScore(ctx context.Context, arg postgres.InsertOrUpdateMlScoreParams) (uuid.UUID, error)
}

type mlScoreService struct {
	mlScoreStorage mlScoreStorage
}

func NewMlScoreService(mlScoreStorage mlScoreStorage) *mlScoreService {
	return &mlScoreService{
		mlScoreStorage: mlScoreStorage,
	}
}

func (s *mlScoreService) InsertOrUpdateMlScore(ctx context.Context, dto dto.CreateMlScoreDTO) (uuid.UUID, error) {
	tracer := otel.Tracer("ml-score-service")
	ctx, span := tracer.Start(ctx, "ml-score-service")
	defer span.End()

	return s.mlScoreStorage.InsertOrUpdateMlScore(ctx, postgres.InsertOrUpdateMlScoreParams{
		ClientID:     uuid.MustParse(dto.ClientID),
		AdvertiserID: uuid.MustParse(dto.AdvertiserID),
		Score:        ads.Normalization(float64(dto.Score)),
	})
}
