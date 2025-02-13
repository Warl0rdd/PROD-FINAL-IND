package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
)

type campaignStorage interface {
	CreateCampaign(ctx context.Context, arg postgres.CreateCampaignParams) (postgres.CreateCampaignRow, error)
	GetCampaignById(ctx context.Context, arg postgres.GetCampaignByIdParams) (entity.Campaign, error)
	GetCampaignWithPagination(ctx context.Context, arg postgres.GetCampaignWithPaginationParams) ([]entity.Campaign, error)
	UpdateCampaign(ctx context.Context, arg postgres.UpdateCampaignParams) (entity.Campaign, error)
	DeleteCampaign(ctx context.Context, arg postgres.DeleteCampaignParams) error
}

type CampaignService struct {
	campaignStorage campaignStorage
}

func NewCampaignService(campaignStorage campaignStorage) *CampaignService {
	return &CampaignService{
		campaignStorage: campaignStorage,
	}
}

// TODO fix gender

func (s *CampaignService) CreateCampaign(ctx context.Context, campaignDTO dto.CreateCampaignDTO) (dto.CampaignDTO, error) {
	// Если нам не передали нижнюю границу возраста и нам подставилось нулевое значение - нас это устраивает
	// А если не передали верхнюю границу возраста - ставим максимальное, что бы сортировка таргетинга по верхней границе возраста не применялась

	var ageTo int32
	if campaignDTO.Targeting.AgeTo != 0 {
		ageTo = campaignDTO.Targeting.AgeTo
	} else {
		ageTo = 999
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
		Targeting:         campaignDTO.Targeting,
	}, nil
}

func (s *CampaignService) GetCampaignById(ctx context.Context, campaignDTO dto.GetCampaignByIDDTO) (dto.CampaignDTO, error) {
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
	}, nil
}

func (s *CampaignService) GetCampaignWithPagination(ctx context.Context, campaignDTO dto.GetCampaignsWithPaginationDTO) ([]dto.CampaignDTO, error) {
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
		})
	}
	return result, nil
}

// TODO проверка на дату старта компании

func (s *CampaignService) UpdateCampaign(ctx context.Context, campaignDTO dto.UpdateCampaignDTO) (dto.CampaignDTO, error) {
	campaign, err := s.campaignStorage.UpdateCampaign(ctx, postgres.UpdateCampaignParams{
		ID:                uuid.MustParse(campaignDTO.CampaignID),
		AdvertiserID:      uuid.MustParse(campaignDTO.AdvertiserID),
		CostPerImpression: campaignDTO.CostPerImpression,
		CostPerClick:      campaignDTO.CostPerClick,
		AdTitle:           campaignDTO.AdTitle,
		AdText:            campaignDTO.AdText,
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
		Gender: entity.CampaignGender(campaignDTO.Targeting.Gender),
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
	}, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignDTO dto.DeleteCampaignDTO) error {
	return s.campaignStorage.DeleteCampaign(ctx, postgres.DeleteCampaignParams{
		ID:           uuid.MustParse(campaignDTO.CampaignID),
		AdvertiserID: uuid.MustParse(campaignDTO.AdvertiserID),
	})
}
