package v1

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type StatsService interface {
	GetDailyStatsByAdvertiserID(ctx context.Context, statsDTO dto.GetStatsByAdvertiserIDDTO) ([]dto.StatsDTO, error)
	GetDailyStatsByCampaignID(ctx context.Context, statsDTO dto.GetStatsByCampaignIDDTO) ([]dto.StatsDTO, error)
	GetStatsByAdvertiserID(ctx context.Context, statsDTO dto.GetStatsByAdvertiserIDDTO) (dto.StatsDTO, error)
	GetStatsByCampaignID(ctx context.Context, statsDTO dto.GetStatsByCampaignIDDTO) (dto.StatsDTO, error)
}

type StatsHandler struct {
	statsService StatsService
	validator    *validator.Validator
}

func NewStatsHandler(app *app.App) *StatsHandler {
	statsStorage := postgres.NewStatsStorage(app.DB)

	return &StatsHandler{
		statsService: service.NewStatsService(statsStorage),
		validator:    app.Validator,
	}
}

func (h *StatsHandler) GetDailyStatsByAdvertiserID(c fiber.Ctx) error {
	var statsDTO dto.GetStatsByAdvertiserIDDTO

	if err := c.Bind().URI(&statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	stats, err := h.statsService.GetDailyStatsByAdvertiserID(c.Context(), statsDTO)
	if err != nil {
		if errors.Is(err, errorz.NotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Stats not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

func (h *StatsHandler) GetDailyStatsByCampaignID(c fiber.Ctx) error {
	var statsDTO dto.GetStatsByCampaignIDDTO

	if err := c.Bind().URI(&statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	stats, err := h.statsService.GetDailyStatsByCampaignID(c.Context(), statsDTO)
	if err != nil {
		if errors.Is(err, errorz.NotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Stats not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

func (h *StatsHandler) GetStatsByAdvertiserID(c fiber.Ctx) error {
	var statsDTO dto.GetStatsByAdvertiserIDDTO

	if err := c.Bind().URI(&statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	stats, err := h.statsService.GetStatsByAdvertiserID(c.Context(), statsDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

func (h *StatsHandler) GetStatsByCampaignID(c fiber.Ctx) error {
	var statsDTO dto.GetStatsByCampaignIDDTO

	if err := c.Bind().URI(&statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(statsDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	stats, err := h.statsService.GetStatsByCampaignID(c.Context(), statsDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

func (h *StatsHandler) Setup(router fiber.Router) {
	router.Get("/stats/advertisers/:advertiserID/campaigns/daily", h.GetDailyStatsByAdvertiserID)
	router.Get("/stats/campaigns/:campaignID/daily", h.GetDailyStatsByCampaignID)
	router.Get("/stats/advertisers/:advertiserID/campaigns", h.GetStatsByAdvertiserID)
	router.Get("/stats/campaigns/:campaignID", h.GetStatsByCampaignID)
}
