package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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
	var MlScoreDTO dto.CreateMlScoreDTO

	if err := c.Bind().Body(&MlScoreDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(MlScoreDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	_, err := h.mlScoreService.InsertOrUpdateMlScore(c.Context(), MlScoreDTO)
	if err != nil {
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
