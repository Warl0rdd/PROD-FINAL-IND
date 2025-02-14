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
	"solution/internal/adapters/logger"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, campaignDTO dto.CreateCampaignDTO) (dto.CampaignDTO, error)
	GetCampaignById(ctx context.Context, campaignDTO dto.GetCampaignByIDDTO) (dto.CampaignDTO, error)
	GetCampaignWithPagination(ctx context.Context, campaignDTO dto.GetCampaignsWithPaginationDTO) ([]dto.CampaignDTO, error)
	UpdateCampaign(ctx context.Context, campaignDTO dto.UpdateCampaignDTO) (dto.CampaignDTO, error)
	DeleteCampaign(ctx context.Context, campaignDTO dto.DeleteCampaignDTO) error
}

type CampaignHandler struct {
	campaignService CampaignService
	validator       *validator.Validator
}

func NewCampaignHandler(app *app.App) *CampaignHandler {
	campaignStorage := postgres.NewCampaignStorage(app.DB)
	dayStorage := redis.NewDayStorage(app.Redis)

	return &CampaignHandler{
		campaignService: service.NewCampaignService(campaignStorage, dayStorage),
		validator:       app.Validator,
	}
}

func (h *CampaignHandler) CreateCampaign(c fiber.Ctx) error {
	var campaignDTO dto.CreateCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := c.Bind().Body(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	created, err := h.campaignService.CreateCampaign(c.Context(), campaignDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *CampaignHandler) GetCampaignById(c fiber.Ctx) error {
	var campaignDTO dto.GetCampaignByIDDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	campaign, err := h.campaignService.GetCampaignById(c.Context(), campaignDTO)
	if err != nil {
		if campaign == (dto.CampaignDTO{}) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Campaign not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(campaign)
}

func (h *CampaignHandler) GetCampaignWithPagination(c fiber.Ctx) error {
	var campaignDTO dto.GetCampaignsWithPaginationDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().Query(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if campaignDTO.Page == 1 {
		campaignDTO.Offset = 0
	} else {
		campaignDTO.Offset = (campaignDTO.Page - 1) * campaignDTO.Limit
	}

	campaigns, err := h.campaignService.GetCampaignWithPagination(c.Context(), campaignDTO)
	if err != nil {
		if len(campaigns) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Campaigns not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(campaigns)
}

func (h *CampaignHandler) UpdateCampaign(c fiber.Ctx) error {

	var campaignDTO dto.UpdateCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().Body(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	updated, err := h.campaignService.UpdateCampaign(c.Context(), campaignDTO)
	if err != nil {
		if errors.Is(err, errorz.Forbidden) {
			return c.Status(fiber.StatusForbidden).JSON(dto.HTTPError{
				Code:    fiber.StatusForbidden,
				Message: "Can't edit clicks and impressions limits after start date",
			})
		}

		if updated == (dto.CampaignDTO{}) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Campaign not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (h *CampaignHandler) DeleteCampaign(c fiber.Ctx) error {
	var campaignDTO dto.DeleteCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err := h.campaignService.DeleteCampaign(c.Context(), campaignDTO)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "Campaign not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *CampaignHandler) Setup(router fiber.Router) {
	group := router.Group("/advertisers")
	group.Post("/:advertiserId/campaigns", h.CreateCampaign)
	group.Get("/:advertiserId/campaigns/:campaignId", h.GetCampaignById)
	group.Get("/:advertiserId/campaigns", h.GetCampaignWithPagination)
	group.Put("/:advertiserId/campaigns/:campaignId", h.UpdateCampaign)
	group.Delete("/:advertiserId/campaigns/:campaignId", h.DeleteCampaign)
}
