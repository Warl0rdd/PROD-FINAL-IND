package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"solution/cmd/app"
	"solution/internal/adapters/LLM"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type LLMService interface {
	GenerateCampaignText(ctx context.Context, generationDTO dto.LLMRequestDTO) ([]string, error)
}

type LLMHandler struct {
	llmService LLMService
	validator  *validator.Validator
}

func NewLLMHandler(app *app.App) *LLMHandler {
	return &LLMHandler{
		llmService: service.NewLLMService(LLM.NewYandexGPTStorage(app.GPT)),
		validator:  app.Validator,
	}
}

func (h *LLMHandler) GenerateCampaignText(c fiber.Ctx) error {
	tracer := otel.Tracer("GenerateCampaignText")
	ctx, span := tracer.Start(c.Context(), "LLMHandler")
	defer span.End()

	var generateDTO dto.LLMRequestDTO

	if err := c.Bind().Body(&generateDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(generateDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	texts, err := h.llmService.GenerateCampaignText(ctx, generateDTO)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	var result []dto.LLMResponseDTO

	for _, v := range texts {
		result = append(result, dto.LLMResponseDTO{
			CampaignText: v,
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *LLMHandler) Setup(router fiber.Router) {
	router.Post("/llm", h.GenerateCampaignText)
}
