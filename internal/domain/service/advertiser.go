package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
)

type advertisersStorage interface {
	CreateAdvertiser(ctx context.Context, arg postgres.CreateAdvertiserParams) (entity.Advertiser, error)
	GetAdvertiserById(ctx context.Context, id uuid.UUID) (entity.Advertiser, error)
}

type advertisersService struct {
	advertisersStorage advertisersStorage
}

func NewAdvertisersService(advertisersStorage advertisersStorage) *advertisersService {
	return &advertisersService{
		advertisersStorage: advertisersStorage,
	}
}

func (s *advertisersService) CreateAdvertiser(ctx context.Context, dto dto.CreateAdvertiserDTO) (entity.Advertiser, error) {
	tracer := otel.Tracer("advertisers-service")
	ctx, span := tracer.Start(ctx, "advertisers-service")
	defer span.End()

	return s.advertisersStorage.CreateAdvertiser(ctx, postgres.CreateAdvertiserParams{
		ID:   uuid.MustParse(dto.AdvertiserID),
		Name: dto.Name,
	})
}

func (s *advertisersService) GetAdvertiserById(ctx context.Context, dto dto.GetAdvertiserByIdDTO) (entity.Advertiser, error) {
	tracer := otel.Tracer("advertisers-service")
	ctx, span := tracer.Start(ctx, "advertisers-service")
	defer span.End()

	return s.advertisersStorage.GetAdvertiserById(ctx, uuid.MustParse(dto.AdvertiserID))
}
