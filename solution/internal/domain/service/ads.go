package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/logger"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	adUtils "solution/internal/domain/utils/ads"
	"solution/internal/domain/utils/learning"
)

type AdsStorage interface {
	GetEligibleAds(ctx context.Context, arg postgres.GetEligibleAdsParams) ([]postgres.GetEligibleAdsRow, error)
	AddImpression(ctx context.Context, arg postgres.AddImpressionParams) error
	AddClick(ctx context.Context, arg postgres.AddClickParams) error
}

type RedisLearningStorage interface {
	GetR0(ctx context.Context) float64
	SetR0(ctx context.Context, r0 float64) error
}

type PostgresLearningStorage interface {
	GetImpressionsForLearning(ctx context.Context) ([]postgres.GetImpressionsForLearningRow, error)
	UpdateLearnedImpression(ctx context.Context, id uuid.UUID) error
}

type adsService struct {
	adsStorage              AdsStorage
	dayStorage              DayStorage
	redisLearningStorage    RedisLearningStorage
	postgresLearningStorage PostgresLearningStorage
}

func NewAdsService(adsStorage AdsStorage, dayStorage DayStorage, learningStorage RedisLearningStorage, postgresLearningStorage PostgresLearningStorage) *adsService {
	return &adsService{
		adsStorage:              adsStorage,
		dayStorage:              dayStorage,
		redisLearningStorage:    learningStorage,
		postgresLearningStorage: postgresLearningStorage,
	}
}

func (s *adsService) GetAds(ctx context.Context, adsDTO dto.GetAdsDTO) (dto.AdDTO, error) {
	tracer := otel.Tracer("ads-service")
	ctx, span := tracer.Start(ctx, "ads-service")
	defer span.End()

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
		score := adUtils.AdScore(ad.CostPerImpression, ad.CostPerClick, float64(ad.Score), s.redisLearningStorage.GetR0(ctx))

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
		ModelScore: maxKey,
	}); errImp != nil {
		return dto.AdDTO{}, errImp
	}

	span.SetAttributes(attribute.Float64("score", maxKey))

	return scores[maxKey], nil
}

func (s *adsService) Click(ctx context.Context, clickDTO dto.AddClickDTO) error {
	tracer := otel.Tracer("ads-service")
	ctx, span := tracer.Start(ctx, "ads-service")
	defer span.End()

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

	go s.AdjustModel()

	return nil
}

func (s *adsService) AdjustModel() {
	tracer := otel.Tracer("ads-service")
	ctx, span := tracer.Start(context.Background(), "ads-service")
	defer span.End()

	oldR0 := s.redisLearningStorage.GetR0(ctx)
	data, err := s.postgresLearningStorage.GetImpressionsForLearning(ctx)
	if err != nil {
		logger.Log.Error(err)
	}

	newR0 := learning.GenNewR0(oldR0, data)

	if err := s.redisLearningStorage.SetR0(ctx, newR0); err != nil {
		logger.Log.Error(err)
	}

	for _, v := range data {
		if err := s.postgresLearningStorage.UpdateLearnedImpression(ctx, v.ID); err != nil {
			logger.Log.Error(err)
		}
	}

	span.SetAttributes(
		attribute.Float64("oldR0", oldR0),
		attribute.Float64("newR0", newR0),
	)
}
