package service

import (
	"context"
	"github.com/google/uuid"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
)

type StatsStorage interface {
	GetDailyStatsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID) ([]postgres.GetDailyStatsByAdvertiserIDRow, error)
	GetDailyStatsByCampaignID(ctx context.Context, id uuid.UUID) ([]postgres.GetDailyStatsByCampaignIDRow, error)
	GetStatsByAdvertiserID(ctx context.Context, advertiserID uuid.UUID) (postgres.GetStatsByAdvertiserIDRow, error)
	GetStatsByCampaignID(ctx context.Context, id uuid.UUID) (postgres.GetStatsByCampaignIDRow, error)
}

type StatsService struct {
	statsStorage StatsStorage
}

func NewStatsService(statsStorage StatsStorage) *StatsService {
	return &StatsService{
		statsStorage: statsStorage,
	}
}

// TODO статы получают цену за клик в конкретный день до/после апдейта
// TODO проверка на параметры, которые можно/нельзя изменять после старта

func (s *StatsService) GetDailyStatsByAdvertiserID(ctx context.Context, statsDTO dto.GetStatsByAdvertiserIDDTO) ([]dto.StatsDTO, error) {
	stats, err := s.statsStorage.GetDailyStatsByAdvertiserID(ctx, uuid.MustParse(statsDTO.AdvertiserID))

	if err != nil {
		return []dto.StatsDTO{}, err
	}

	if len(stats) == 0 {
		return []dto.StatsDTO{}, errorz.NotFound
	}

	var statsDTOs []dto.StatsDTO
	for _, stat := range stats {
		statsDTOs = append(statsDTOs, dto.StatsDTO{
			ImpressionsCount: int(stat.ImpressionsCount),
			ClicksCount:      int(stat.ClicksCount),
			Conversion:       stat.Conversion,
			SpentImpressions: stat.SpentImpressions,
			SpentClicks:      stat.SpentClicks,
			SpentTotal:       stat.SpentTotal,
			Day:              int(stat.Day),
		})
	}

	return statsDTOs, nil
}

func (s *StatsService) GetDailyStatsByCampaignID(ctx context.Context, statsDTO dto.GetStatsByCampaignIDDTO) ([]dto.StatsDTO, error) {
	stats, err := s.statsStorage.GetDailyStatsByCampaignID(ctx, uuid.MustParse(statsDTO.CampaignID))

	if err != nil {
		return []dto.StatsDTO{}, err
	}

	if len(stats) == 0 {
		return []dto.StatsDTO{}, errorz.NotFound
	}

	var statsDTOs []dto.StatsDTO
	for _, stat := range stats {
		statsDTOs = append(statsDTOs, dto.StatsDTO{
			ImpressionsCount: int(stat.ImpressionsCount),
			ClicksCount:      int(stat.ClicksCount),
			Conversion:       stat.Conversion,
			SpentImpressions: stat.SpentImpressions,
			SpentClicks:      stat.SpentClicks,
			SpentTotal:       stat.SpentTotal,
			Day:              int(stat.Day),
		})
	}

	return statsDTOs, nil
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
