package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
)

type ModerationStorage interface {
	GetCampaignsForModeration(ctx context.Context, arg postgres.GetCampaignsForModerationParams) ([]postgres.GetCampaignsForModerationRow, error)
	Approve(ctx context.Context, id uuid.UUID) error
	Reject(ctx context.Context, id uuid.UUID) error
}

type moderationService struct {
	moderationStorage ModerationStorage
}

func NewModerationService(moderationStorage ModerationStorage) *moderationService {
	return &moderationService{
		moderationStorage: moderationStorage,
	}
}

func (s *moderationService) GetCampaignsForModeration(ctx context.Context, moderationDTO dto.GetCampaignsForModerationDTO) ([]dto.CampaignForModerationDTO, error) {
	tracer := otel.Tracer("moderation-service")
	ctx, span := tracer.Start(ctx, "moderation-service")
	defer span.End()

	campaigns, err := s.moderationStorage.GetCampaignsForModeration(ctx, postgres.GetCampaignsForModerationParams{
		Limit:  int32(moderationDTO.Limit),
		Offset: int32(moderationDTO.Offset),
	})
	if err != nil {
		return nil, err
	}
	var result []dto.CampaignForModerationDTO
	for _, campaign := range campaigns {
		result = append(result, dto.CampaignForModerationDTO{
			ID:           campaign.ID.String(),
			AdvertiserID: campaign.AdvertiserID.String(),
			AdTitle:      campaign.AdTitle,
			AdText:       campaign.AdText,
		})
	}
	return result, nil
}

func (s *moderationService) Approve(ctx context.Context, id string) error {
	tracer := otel.Tracer("moderation-service")
	ctx, span := tracer.Start(ctx, "moderation-service")
	defer span.End()

	return s.moderationStorage.Approve(ctx, uuid.MustParse(id))
}

func (s *moderationService) Reject(ctx context.Context, id string) error {
	tracer := otel.Tracer("moderation-service")
	ctx, span := tracer.Start(ctx, "moderation-service")
	defer span.End()

	return s.moderationStorage.Reject(ctx, uuid.MustParse(id))
}
