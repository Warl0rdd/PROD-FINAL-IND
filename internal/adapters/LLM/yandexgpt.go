package LLM

import (
	"context"
	"fmt"
	"github.com/sheeiavellie/go-yandexgpt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"os"
	"solution/internal/adapters/logger"
)

type yandexGPTStorage struct {
	client *yandexgpt.YandexGPTClient
}

func NewYandexGPTStorage(client *yandexgpt.YandexGPTClient) *yandexGPTStorage {
	return &yandexGPTStorage{
		client: client,
	}
}

const systemPrompt = "Ты - модель для генерации текстов к рекламных кампаниям. Они должны быть продающими, краткими и лаконичными, буквально пара предложений. В промпте тебе будет переданы два параметра - заголовок рекламной кампании (title) и название рекламодателя (advertiser_name). Твоя задача - сгенерировать текст кампании. Нужен только он, без каких либо дополнительных сообщений или текста. Старайся не повторять словосочетания из заголовка и больше раскрывать тему."

func (s *yandexGPTStorage) GenerateCampaignText(ctx context.Context, campaignTitle, advertiserName string) ([]yandexgpt.YandexGPTAlternative, error) {
	tracer := otel.Tracer("GenerateCampaignText")
	ctx, span := tracer.Start(ctx, "YandexGPTStorage")
	defer span.End()

	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(os.Getenv("YANDEX_GPT_CATALOG_ID"), yandexgpt.YandexGPT4ModelLite),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			MaxTokens:   500,
			Stream:      false,
			Temperature: 0.7,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{

				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: systemPrompt,
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: fmt.Sprintf("title: %s\n advertiser_name: %s", campaignTitle, advertiserName),
			},
		},
	}

	response, err := s.client.GetCompletion(ctx, request)
	if err != nil {
		span.RecordError(err)
		logger.Log.Errorf("Failed to generate campaign text: %v", err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("text", response.Result.Alternatives[0].Message.Text),
		attribute.String("input-tokens", response.Result.Usage.InputTokens),
		attribute.String("completion-tokens", response.Result.Usage.CompletionTokens),
		attribute.String("total-tokens", response.Result.Usage.TotalTokens),
	)
	return response.Result.Alternatives, nil

}
