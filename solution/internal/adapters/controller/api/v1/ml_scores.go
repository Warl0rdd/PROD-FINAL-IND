package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type MLScoreService interface {
	InsertOrUpdateMlScore(ctx context.Context, dto dto.CreateMlScoreDTO) (uuid.UUID, error)
}

type MlScoreHandler struct {
	mlScoreService MLScoreService
	validator      *validator.Validator
}

func NewMlScoreHandler(app *app.App) *MlScoreHandler {
	mlScoreStorage := postgres.NewMlScoreStorage(app.DB)

	return &MlScoreHandler{
		mlScoreService: service.NewMlScoreService(mlScoreStorage),
		validator:      app.Validator,
	}
}

func (h *MlScoreHandler) InsertOrUpdateMlScore(c fiber.Ctx) error {
	tracer := otel.Tracer("ml-score-handler")
	ctx, span := tracer.Start(c.Context(), "InsertOrUpdateMlScore")
	defer span.End()

	var MlScoreDTO dto.CreateMlScoreDTO

	if err := c.Bind().Body(&MlScoreDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(MlScoreDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.String("advertiserId", MlScoreDTO.AdvertiserID),
		attribute.String("clientId", MlScoreDTO.ClientID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/ml-scores"),
	)

	_, err := h.mlScoreService.InsertOrUpdateMlScore(ctx, MlScoreDTO)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

func (h *MlScoreHandler) Setup(router fiber.Router) {
	router.Post("/ml-scores", h.InsertOrUpdateMlScore)
}
