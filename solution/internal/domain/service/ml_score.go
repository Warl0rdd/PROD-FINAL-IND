package service

import (
	"context"
	"github.com/google/uuid"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
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
	return s.mlScoreStorage.InsertOrUpdateMlScore(ctx, postgres.InsertOrUpdateMlScoreParams{
		ClientID:     uuid.MustParse(dto.ClientID),
		AdvertiserID: uuid.MustParse(dto.AdvertiserID),
		Score:        int32(dto.Score),
	})
}
