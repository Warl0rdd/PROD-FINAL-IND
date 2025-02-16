package service

import (
	"context"
	"github.com/sheeiavellie/go-yandexgpt"
	"go.opentelemetry.io/otel"
	"solution/internal/domain/dto"
)

type LLMStorage interface {
	GenerateCampaignText(ctx context.Context, campaignTitle, advertiserName string) ([]yandexgpt.YandexGPTAlternative, error)
}

type LLMService struct {
	llmStorage LLMStorage
}

func NewLLMService(llmStorage LLMStorage) *LLMService {
	return &LLMService{
		llmStorage: llmStorage,
	}
}

func (s *LLMService) GenerateCampaignText(ctx context.Context, generationDTO dto.LLMRequestDTO) ([]string, error) {
	tracer := otel.Tracer("GenerateCampaignText")
	ctx, span := tracer.Start(ctx, "LLMService")
	defer span.End()

	alternatives, err := s.llmStorage.GenerateCampaignText(ctx, generationDTO.CampaignTitle, generationDTO.AdvertiserName)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var result []string
	for _, alternative := range alternatives {
		result = append(result, alternative.Message.Text)
	}
	return result, nil
}
