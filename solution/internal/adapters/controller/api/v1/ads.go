package v1

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/database/redis"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type AdsService interface {
	GetAds(ctx context.Context, adsDTO dto.GetAdsDTO) (dto.AdDTO, error)
	Click(ctx context.Context, clickDTO dto.AddClickDTO) error
}

type AdsHandler struct {
	adsService AdsService
	validator  *validator.Validator
}

func NewAdsHandler(app *app.App) *AdsHandler {
	adsStorage := postgres.NewAdsStorage(app.DB)
	dayStorage := redis.NewDayStorage(app.Redis)

	redisLearningStorage := redis.NewLearningStorage(app.Redis)
	postgresLearningStorage := postgres.NewLearningStorage(app.DB)

	return &AdsHandler{
		adsService: service.NewAdsService(adsStorage, dayStorage, redisLearningStorage, postgresLearningStorage),
		validator:  app.Validator,
	}
}

func (h *AdsHandler) GetAds(c fiber.Ctx) error {
	var adsDTO dto.GetAdsDTO

	if err := c.Bind().Query(&adsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(adsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	ads, err := h.adsService.GetAds(c.Context(), adsDTO)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Ads not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ads)
}

func (h *AdsHandler) Click(c fiber.Ctx) error {
	var clickDTO dto.AddClickDTO

	if err := c.Bind().Body(&clickDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().URI(&clickDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(clickDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err := h.adsService.Click(c.Context(), clickDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *AdsHandler) Setup(router fiber.Router) {
	router.Get("/ads", h.GetAds)
	router.Post("/ads/:adID/click", h.Click)
}
