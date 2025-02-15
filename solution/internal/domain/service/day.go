package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"solution/internal/domain/dto"
)

type DayStorage interface {
	SetDay(ctx context.Context, day int) error
	GetDay(ctx context.Context) (int, error)
}

type dayService struct {
	dayStorage DayStorage
}

func NewDayService(dayStorage DayStorage) *dayService {
	return &dayService{
		dayStorage: dayStorage,
	}
}

func (s *dayService) SetDay(ctx context.Context, dto dto.SetDayDTO) (int, error) {
	tracer := otel.Tracer("day-service")
	ctx, span := tracer.Start(ctx, "day-service")
	defer span.End()

	day := dto.CurrentDate

	if day == 0 {
		currentDay, err := s.dayStorage.GetDay(ctx)
		if err != nil {
			return 0, err
		}
		day = currentDay + 1
	}

	return day, s.dayStorage.SetDay(ctx, day)
}

func (s *dayService) GetDay(ctx context.Context) (int, error) {
	return s.dayStorage.GetDay(ctx)
}
