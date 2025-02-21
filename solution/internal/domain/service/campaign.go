package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
)

type campaignStorage interface {
	CreateCampaign(ctx context.Context, arg postgres.CreateCampaignParams) (postgres.CreateCampaignRow, error)
	GetCampaignById(ctx context.Context, arg postgres.GetCampaignByIdParams) (entity.Campaign, error)
	GetCampaignByIdInsecure(ctx context.Context, campaignId uuid.UUID) (entity.Campaign, error)
	GetCampaignWithPagination(ctx context.Context, arg postgres.GetCampaignWithPaginationParams) ([]entity.Campaign, error)
	UpdateCampaign(ctx context.Context, arg postgres.UpdateCampaignParams) (entity.Campaign, error)
	DeleteCampaign(ctx context.Context, arg postgres.DeleteCampaignParams) error
}

type CampaignService struct {
	campaignStorage campaignStorage
	dayStorage      DayStorage
}

func NewCampaignService(campaignStorage campaignStorage, dayStorage DayStorage) *CampaignService {
	return &CampaignService{
		campaignStorage: campaignStorage,
		dayStorage:      dayStorage,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, campaignDTO dto.CreateCampaignDTO) (dto.CampaignDTO, error) {
	// Если нам не передали нижнюю границу возраста и нам подставилось нулевое значение - нас это устраивает
	// А если не передали верхнюю границу возраста - ставим максимальное, что бы сортировка таргетинга по верхней границе возраста не применялась

	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	var ageTo int32
	if campaignDTO.Targeting.AgeTo != 0 {
		ageTo = campaignDTO.Targeting.AgeTo
	} else {
		ageTo = 999
	}

	if campaignDTO.Targeting.AgeFrom != 0 && campaignDTO.Targeting.AgeTo != 0 && campaignDTO.Targeting.AgeFrom > campaignDTO.Targeting.AgeTo {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.StartDate > campaignDTO.EndDate {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.Gender != "ALL" && campaignDTO.Targeting.Gender != "MALE" && campaignDTO.Targeting.Gender != "FEMALE" && campaignDTO.Targeting.Gender != "" {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.Gender == "" {
		campaignDTO.Targeting.Gender = "ALL"
	}

	created, err := s.campaignStorage.CreateCampaign(ctx, postgres.CreateCampaignParams{
		AdvertiserID:      uuid.MustParse(campaignDTO.AdvertiserID),
		ImpressionLimit:   campaignDTO.ImpressionsLimit,
		ClicksLimit:       campaignDTO.ClicksLimit,
		CostPerImpression: campaignDTO.CostPerImpression,
		CostPerClick:      campaignDTO.CostPerClick,
		AdTitle:           campaignDTO.AdTitle,
		AdText:            campaignDTO.AdText,
		StartDate:         campaignDTO.StartDate,
		EndDate:           campaignDTO.EndDate,
		AgeFrom: pgtype.Int4{
			Int32: campaignDTO.Targeting.AgeFrom,
			Valid: true,
		},
		AgeTo: pgtype.Int4{
			Int32: ageTo,
			Valid: true,
		},
		Location: pgtype.Text{
			String: campaignDTO.Targeting.Location,
			Valid:  true,
		},
		Gender: entity.CampaignGender(campaignDTO.Targeting.Gender),
	})

	if err != nil {
		return dto.CampaignDTO{}, err
	}
	return dto.CampaignDTO{
		CampaignID:        created.ID.String(),
		AdvertiserID:      created.AdvertiserID.String(),
		ImpressionsLimit:  campaignDTO.ImpressionsLimit,
		ClicksLimit:       campaignDTO.ClicksLimit,
		CostPerImpression: campaignDTO.CostPerImpression,
		CostPerClick:      campaignDTO.CostPerClick,
		AdTitle:           campaignDTO.AdTitle,
		AdText:            campaignDTO.AdText,
		StartDate:         campaignDTO.StartDate,
		EndDate:           campaignDTO.EndDate,
		Targeting: dto.Target{
			AgeFrom:  campaignDTO.Targeting.AgeFrom,
			AgeTo:    ageTo,
			Location: campaignDTO.Targeting.Location,
			Gender:   campaignDTO.Targeting.Gender,
		},
	}, nil
}

func (s *CampaignService) GetCampaignById(ctx context.Context, campaignDTO dto.GetCampaignByIDDTO) (dto.CampaignDTO, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	campaign, err := s.campaignStorage.GetCampaignById(ctx, postgres.GetCampaignByIdParams{
		AdvertiserID: uuid.MustParse(campaignDTO.AdvertiserID),
		ID:           uuid.MustParse(campaignDTO.CampaignID),
	})
	if err != nil {
		return dto.CampaignDTO{}, err
	}
	return dto.CampaignDTO{
		CampaignID:        campaign.ID.String(),
		AdvertiserID:      campaign.AdvertiserID.String(),
		ImpressionsLimit:  campaign.ImpressionLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           campaign.AdTitle,
		AdText:            campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: dto.Target{
			AgeFrom:  campaign.AgeFrom.Int32,
			AgeTo:    campaign.AgeTo.Int32,
			Location: campaign.Location.String,
			Gender:   string(campaign.Gender),
		},
		Approved: campaign.Approved,
	}, nil
}

func (s *CampaignService) GetCampaignWithPagination(ctx context.Context, campaignDTO dto.GetCampaignsWithPaginationDTO) ([]dto.CampaignDTO, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	campaigns, err := s.campaignStorage.GetCampaignWithPagination(ctx, postgres.GetCampaignWithPaginationParams{
		AdvertiserID: uuid.MustParse(campaignDTO.AdvertiserID),
		Limit:        campaignDTO.Limit,
		Offset:       campaignDTO.Offset,
	})
	if err != nil {
		return nil, err
	}
	var result []dto.CampaignDTO
	for _, campaign := range campaigns {
		result = append(result, dto.CampaignDTO{
			CampaignID:        campaign.ID.String(),
			AdvertiserID:      campaign.AdvertiserID.String(),
			ImpressionsLimit:  campaign.ImpressionLimit,
			ClicksLimit:       campaign.ClicksLimit,
			CostPerImpression: campaign.CostPerImpression,
			CostPerClick:      campaign.CostPerClick,
			AdTitle:           campaign.AdTitle,
			AdText:            campaign.AdText,
			StartDate:         campaign.StartDate,
			EndDate:           campaign.EndDate,
			Targeting: dto.Target{
				AgeFrom:  campaign.AgeFrom.Int32,
				AgeTo:    campaign.AgeTo.Int32,
				Location: campaign.Location.String,
				Gender:   string(campaign.Gender),
			},
			Approved: campaign.Approved,
		})
	}
	return result, nil
}

func (s *CampaignService) GetCampaignByIdInsecure(ctx context.Context, campaignId string) (dto.CampaignDTO, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	campaign, err := s.campaignStorage.GetCampaignByIdInsecure(ctx, uuid.MustParse(campaignId))
	if err != nil {
		return dto.CampaignDTO{}, err
	}

	return dto.CampaignDTO{
		CampaignID:        campaign.ID.String(),
		AdvertiserID:      campaign.AdvertiserID.String(),
		ImpressionsLimit:  campaign.ImpressionLimit,
		ClicksLimit:       campaign.ClicksLimit,
		CostPerImpression: campaign.CostPerImpression,
		CostPerClick:      campaign.CostPerClick,
		AdTitle:           campaign.AdTitle,
		AdText:            campaign.AdText,
		StartDate:         campaign.StartDate,
		EndDate:           campaign.EndDate,
		Targeting: dto.Target{
			AgeFrom:  campaign.AgeFrom.Int32,
			AgeTo:    campaign.AgeTo.Int32,
			Location: campaign.Location.String,
			Gender:   string(campaign.Gender),
		},
		Approved: campaign.Approved,
	}, nil
}

func (s *CampaignService) UpdateCampaign(ctx context.Context, campaignDTO dto.UpdateCampaignDTO) (dto.CampaignDTO, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	day, err := s.dayStorage.GetDay(ctx)

	if err != nil {
		return dto.CampaignDTO{}, err
	}

	span.SetAttributes(attribute.Int("day", day))

	campaign, err := s.campaignStorage.GetCampaignById(ctx, postgres.GetCampaignByIdParams{
		AdvertiserID: uuid.MustParse(campaignDTO.AdvertiserID),
		ID:           uuid.MustParse(campaignDTO.CampaignID),
	})

	if err != nil || campaign == (entity.Campaign{}) {
		return dto.CampaignDTO{}, err
	}

	if campaign.StartDate <= int32(day) {
		if campaignDTO.ClicksLimit != nil || campaignDTO.ImpressionsLimit != nil || campaignDTO.StartDate != nil || campaignDTO.EndDate != nil || campaignDTO.Targeting.Gender != nil ||
			campaignDTO.Targeting.AgeFrom != 0 || campaignDTO.Targeting.AgeTo != 0 || campaignDTO.Targeting.Location != "" {
			return dto.CampaignDTO{}, errorz.BadRequest
		}
	}

	if campaignDTO.StartDate != nil && (*campaignDTO.StartDate < day || *campaignDTO.StartDate > int(campaign.EndDate)) {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.EndDate != nil && (*campaignDTO.EndDate < day || *campaignDTO.EndDate < int(campaign.StartDate)) {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.StartDate != nil && campaignDTO.EndDate != nil && *campaignDTO.StartDate > *campaignDTO.EndDate {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.Gender != nil && (*campaignDTO.Targeting.Gender != "MALE" && *campaignDTO.Targeting.Gender != "FEMALE" && *campaignDTO.Targeting.Gender != "ALL") {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.AgeFrom != 0 && campaignDTO.Targeting.AgeTo != 0 && campaignDTO.Targeting.AgeFrom > campaignDTO.Targeting.AgeTo {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.AgeFrom != 0 && campaignDTO.Targeting.AgeFrom > campaign.AgeTo.Int32 {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	if campaignDTO.Targeting.AgeTo != 0 && campaignDTO.Targeting.AgeTo < campaign.AgeFrom.Int32 {
		return dto.CampaignDTO{}, errorz.BadRequest
	}

	updatedCampaign, err := s.campaignStorage.UpdateCampaign(ctx, postgres.UpdateCampaignParams{
		ID:                uuid.MustParse(campaignDTO.CampaignID),
		AdvertiserID:      uuid.MustParse(campaignDTO.AdvertiserID),
		CostPerImpression: campaignDTO.CostPerImpression,
		CostPerClick:      campaignDTO.CostPerClick,
		AdTitle:           campaignDTO.AdTitle,
		AdText:            campaignDTO.AdText,
		ImpressionLimit:   campaignDTO.ImpressionsLimit,
		ClicksLimit:       campaignDTO.ClicksLimit,
		StartDate:         campaignDTO.StartDate,
		EndDate:           campaignDTO.EndDate,
		AgeFrom: pgtype.Int4{
			Int32: campaignDTO.Targeting.AgeFrom,
			Valid: true,
		},
		AgeTo: pgtype.Int4{
			Int32: campaignDTO.Targeting.AgeTo,
			Valid: true,
		},
		Location: pgtype.Text{
			String: campaignDTO.Targeting.Location,
			Valid:  true,
		},
		Gender: campaignDTO.Targeting.Gender,
	})

	if err != nil {
		return dto.CampaignDTO{}, err
	}
	return dto.CampaignDTO{
		CampaignID:        updatedCampaign.ID.String(),
		AdvertiserID:      updatedCampaign.AdvertiserID.String(),
		ImpressionsLimit:  updatedCampaign.ImpressionLimit,
		ClicksLimit:       updatedCampaign.ClicksLimit,
		CostPerImpression: updatedCampaign.CostPerImpression,
		CostPerClick:      updatedCampaign.CostPerClick,
		AdTitle:           updatedCampaign.AdTitle,
		AdText:            updatedCampaign.AdText,
		StartDate:         updatedCampaign.StartDate,
		EndDate:           updatedCampaign.EndDate,
		Targeting: dto.Target{
			AgeFrom:  updatedCampaign.AgeFrom.Int32,
			AgeTo:    updatedCampaign.AgeTo.Int32,
			Location: updatedCampaign.Location.String,
			Gender:   string(updatedCampaign.Gender),
		},
		Approved: updatedCampaign.Approved,
	}, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignDTO dto.DeleteCampaignDTO) error {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	return s.campaignStorage.DeleteCampaign(ctx, postgres.DeleteCampaignParams{
		ID:           uuid.MustParse(campaignDTO.CampaignID),
		AdvertiserID: uuid.MustParse(campaignDTO.AdvertiserID),
	})
}
