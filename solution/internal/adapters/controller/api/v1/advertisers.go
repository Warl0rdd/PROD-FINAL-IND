package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/logger"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
	"solution/internal/domain/service"
)

type AdvertiserService interface {
	CreateAdvertiser(ctx context.Context, dto dto.CreateAdvertiserDTO) (entity.Advertiser, error)
	GetAdvertiserById(ctx context.Context, dto dto.GetAdvertiserByIdDTO) (entity.Advertiser, error)
}

type AdvertiserHandler struct {
	advertiserService AdvertiserService
	validator         *validator.Validator
}

func NewAdvertiserHandler(app *app.App) *AdvertiserHandler {
	advertiserStorage := postgres.NewAdvertiserStorage(app.DB)

	return &AdvertiserHandler{
		advertiserService: service.NewAdvertisersService(advertiserStorage),
		validator:         app.Validator,
	}
}

func (h *AdvertiserHandler) CreateAdvertiser(c fiber.Ctx) error {
	var DTOs []dto.CreateAdvertiserDTO

	if err := c.Bind().Body(&DTOs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("DTOs: %v", DTOs)

	for _, v := range DTOs {
		if err := h.validator.ValidateData(v); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	_, err := h.advertiserService.CreateAdvertiser(c.Context(), DTOs[0])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(DTOs)
}

func (h *AdvertiserHandler) GetAdvertiserById(c fiber.Ctx) error {
	var getAdvertiserDTO dto.GetAdvertiserByIdDTO

	if err := c.Bind().URI(&getAdvertiserDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(getAdvertiserDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	advertiser, err := h.advertiserService.GetAdvertiserById(c.Context(), getAdvertiserDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.CreateAdvertiserDTO{
		AdvertiserID: advertiser.ID.String(),
		Name:         advertiser.Name,
	})
}

func (h *AdvertiserHandler) Setup(router fiber.Router) {
	advertiserGroup := router.Group("/advertisers")
	advertiserGroup.Post("/bulk", h.CreateAdvertiser)
	advertiserGroup.Get("/:advertiserId", h.GetAdvertiserById)
}
