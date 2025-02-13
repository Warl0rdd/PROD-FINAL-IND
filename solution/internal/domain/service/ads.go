package service

import (
	"context"
	"github.com/google/uuid"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/logger"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	adUtils "solution/internal/domain/utils/ads"
)

type AdsStorage interface {
	GetEligibleAds(ctx context.Context, arg postgres.GetEligibleAdsParams) ([]postgres.GetEligibleAdsRow, error)
	AddImpression(ctx context.Context, arg postgres.AddImpressionParams) error
	AddClick(ctx context.Context, arg postgres.AddClickParams) error
}

type adsService struct {
	adsStorage AdsStorage
	dayStorage DayStorage
}

func NewAdsService(adsStorage AdsStorage, dayStorage DayStorage) *adsService {
	return &adsService{
		adsStorage: adsStorage,
		dayStorage: dayStorage,
	}
}

func (s *adsService) GetAds(ctx context.Context, adsDTO dto.GetAdsDTO) (dto.AdDTO, error) {
	day, err := s.dayStorage.GetDay(ctx)
	if err != nil {
		return dto.AdDTO{}, err
	}

	ads, errGet := s.adsStorage.GetEligibleAds(ctx, postgres.GetEligibleAdsParams{
		ClientID: uuid.MustParse(adsDTO.ClientID),
		Day:      int32(day),
	})

	if errGet != nil {
		return dto.AdDTO{}, errGet
	}

	if len(ads) == 0 {
		return dto.AdDTO{}, errorz.NotFound
	}

	scores := make(map[float64]dto.AdDTO, len(ads))

	for _, ad := range ads {
		score := adUtils.AdScore(ad.CostPerImpression, ad.CostPerClick, float64(ad.Score))

		scores[score] = dto.AdDTO{
			AdID:         ad.ID.String(),
			AdvertiserID: ad.AdvertiserID.String(),
			AdTitle:      ad.AdTitle,
			AdText:       ad.AdText,
		}

		logger.Log.Debugf("ad: %v, score: %f", ad, score)
	}

	maxKey := 0.0
	for key := range scores {
		if key > maxKey {
			maxKey = key
		}
	}

	if maxKey == 0.0 {
		return dto.AdDTO{}, errorz.NotFound
	}

	if errImp := s.adsStorage.AddImpression(ctx, postgres.AddImpressionParams{
		ClientID:   uuid.MustParse(adsDTO.ClientID),
		CampaignID: uuid.MustParse(scores[maxKey].AdID),
		Day:        int32(day),
	}); errImp != nil {
		return dto.AdDTO{}, errImp
	}

	return scores[maxKey], nil
}

func (s *adsService) Click(ctx context.Context, clickDTO dto.AddClickDTO) error {
	day, err := s.dayStorage.GetDay(ctx)
	if err != nil {
		return err
	}

	if errClick := s.adsStorage.AddClick(ctx, postgres.AddClickParams{
		ClientID:   uuid.MustParse(clickDTO.ClientID),
		CampaignID: uuid.MustParse(clickDTO.AdID),
		Day:        int32(day),
	}); errClick != nil {
		return errClick
	}

	return nil
}
