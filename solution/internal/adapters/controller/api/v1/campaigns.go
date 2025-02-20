package v1

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	GetCampaignByIdInsecure(ctx context.Context, campaignId string) (dto.CampaignDTO, error)
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
	tracer := otel.Tracer("campaign-handler")
	ctx, span := tracer.Start(c.Context(), "CreateCampaign")
	defer span.End()

	var campaignDTO dto.CreateCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := c.Bind().Body(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if campaignDTO.StartDate > campaignDTO.EndDate {
		span.RecordError(errorz.BadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "start date must be less than end date",
		})
	}

	span.SetAttributes(
		attribute.String("advertiserId", campaignDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/campaigns"),
	)

	created, err := h.campaignService.CreateCampaign(ctx, campaignDTO)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *CampaignHandler) GetCampaignById(c fiber.Ctx) error {
	tracer := otel.Tracer("campaign-handler")
	ctx, span := tracer.Start(c.Context(), "GetCampaignById")
	defer span.End()

	var campaignDTO dto.GetCampaignByIDDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	logger.Log.Debugf("campaignDTO: %v", campaignDTO)

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	campaign, err := h.campaignService.GetCampaignById(ctx, campaignDTO)

	span.SetAttributes(
		attribute.String("advertiserId", campaignDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/campaigns/{campaignId}"),
	)

	if err != nil {
		span.RecordError(err)
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
	tracer := otel.Tracer("campaign-handler")
	ctx, span := tracer.Start(c.Context(), "GetCampaignWithPagination")
	defer span.End()

	var campaignDTO dto.GetCampaignsWithPaginationDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().Query(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		span.RecordError(err)
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

	span.SetAttributes(
		attribute.String("advertiserId", campaignDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/campaigns"),
	)

	campaigns, err := h.campaignService.GetCampaignWithPagination(ctx, campaignDTO)
	if err != nil {
		span.RecordError(err)
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
	tracer := otel.Tracer("campaign-handler")
	ctx, span := tracer.Start(c.Context(), "UpdateCampaign")
	defer span.End()

	var campaignDTO dto.UpdateCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().Body(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.String("advertiserId", campaignDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/campaigns/{campaignId}"),
	)

	updated, err := h.campaignService.UpdateCampaign(ctx, campaignDTO)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, errorz.Forbidden) {
			return c.Status(fiber.StatusForbidden).JSON(dto.HTTPError{
				Code:    fiber.StatusForbidden,
				Message: "Can't edit clicks and impressions limits after start date",
			})
		} else if errors.Is(err, errorz.BadRequest) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
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
	tracer := otel.Tracer("campaign-handler")
	ctx, span := tracer.Start(c.Context(), "DeleteCampaign")
	defer span.End()

	var campaignDTO dto.DeleteCampaignDTO

	if err := c.Bind().URI(&campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(campaignDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if campaign, err := h.campaignService.GetCampaignById(ctx, dto.GetCampaignByIDDTO{
		CampaignID:   campaignDTO.CampaignID,
		AdvertiserID: campaignDTO.AdvertiserID,
	}); err != nil || campaign == (dto.CampaignDTO{}) {
		span.RecordError(err)
		return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: "Campaign not found",
		})
	}

	span.SetAttributes(
		attribute.String("advertiserId", campaignDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}/campaigns/{campaignId}"),
	)

	err := h.campaignService.DeleteCampaign(ctx, campaignDTO)
	if err != nil {
		span.RecordError(err)

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
