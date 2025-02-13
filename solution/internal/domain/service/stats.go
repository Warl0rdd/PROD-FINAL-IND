package service

import (
	"context"
	"github.com/google/uuid"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
)

type StatsStorage interface {
	GetDailyStatsByAdvertiserID(ctx context.Context, arg postgres.GetDailyStatsByAdvertiserIDParams) (postgres.GetDailyStatsByAdvertiserIDRow, error)
	GetDailyStatsByCampaignID(ctx context.Context, arg postgres.GetDailyStatsByCampaignIDParams) (postgres.GetDailyStatsByCampaignIDRow, error)
	GetStatsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID) (postgres.GetStatsByAdvertiserIDRow, error)
	GetStatsByCampaignID(ctx context.Context, id uuid.UUID) (postgres.GetStatsByCampaignIDRow, error)
}

type StatsService struct {
	statsStorage StatsStorage
	dayStorage   DayStorage
}

func NewStatsService(statsStorage StatsStorage, dayStorage DayStorage) *StatsService {
	return &StatsService{
		statsStorage: statsStorage,
		dayStorage:   dayStorage,
	}
}

func (s *StatsService) GetDailyStatsByAdvertiserID(ctx context.Context, statsDTO dto.GetStatsByAdvertiserIDDTO) (dto.StatsDTO, error) {
	day, err := s.dayStorage.GetDay(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}

	arg := postgres.GetDailyStatsByAdvertiserIDParams{
		AdvertiserID: uuid.MustParse(statsDTO.AdvertiserID),
		Day:          int32(day),
	}

	stats, err := s.statsStorage.GetDailyStatsByAdvertiserID(ctx, arg)

	if err != nil {
		return dto.StatsDTO{}, err
	}

	return dto.StatsDTO{
		ImpressionsCount: int(stats.ImpressionsCount),
		ClicksCount:      int(stats.ClicksCount),
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}, nil
}

func (s *StatsService) GetDailyStatsByCampaignID(ctx context.Context, statsDTO dto.GetStatsByCampaignIDDTO) (dto.StatsDTO, error) {
	day, err := s.dayStorage.GetDay(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}

	arg := postgres.GetDailyStatsByCampaignIDParams{
		ID:  uuid.MustParse(statsDTO.CampaignID),
		Day: int32(day),
	}

	stats, err := s.statsStorage.GetDailyStatsByCampaignID(ctx, arg)

	if err != nil {
		return dto.StatsDTO{}, err
	}

	return dto.StatsDTO{
		ImpressionsCount: int(stats.ImpressionsCount),
		ClicksCount:      int(stats.ClicksCount),
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}, nil
}

func (s *StatsService) GetStatsByAdvertiserID(ctx context.Context, statsDTO dto.GetStatsByAdvertiserIDDTO) (dto.StatsDTO, error) {
	stats, err := s.statsStorage.GetStatsByAdvertiserID(ctx, uuid.MustParse(statsDTO.AdvertiserID))

	if err != nil {
		return dto.StatsDTO{}, err
	}

	return dto.StatsDTO{
		ImpressionsCount: int(stats.TotalImpressions),
		ClicksCount:      int(stats.TotalClicks),
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}, nil
}

func (s *StatsService) GetStatsByCampaignID(ctx context.Context, statsDTO dto.GetStatsByCampaignIDDTO) (dto.StatsDTO, error) {
	stats, err := s.statsStorage.GetStatsByCampaignID(ctx, uuid.MustParse(statsDTO.CampaignID))

	if err != nil {
		return dto.StatsDTO{}, err
	}

	return dto.StatsDTO{
		ImpressionsCount: int(stats.ImpressionsCount),
		ClicksCount:      int(stats.ClicksCount),
		Conversion:       stats.Conversion,
		SpentImpressions: stats.SpentImpressions,
		SpentClicks:      stats.SpentClicks,
		SpentTotal:       stats.SpentTotal,
	}, nil
}
