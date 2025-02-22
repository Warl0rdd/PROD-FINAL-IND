package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/database/redis"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type ModerationService interface {
	GetCampaignsForModeration(ctx context.Context, moderationDTO dto.GetCampaignsForModerationDTO) ([]dto.CampaignForModerationDTO, error)
	Approve(ctx context.Context, id string) error
	Reject(ctx context.Context, id string) error
}

type ModerationHandler struct {
	ModerationService ModerationService
	CampaignService   CampaignService
	validator         *validator.Validator
}

func NewModerationHandler(app *app.App) *ModerationHandler {
	return &ModerationHandler{
		ModerationService: service.NewModerationService(postgres.NewModerationStorage(app.DB)),
		CampaignService:   service.NewCampaignService(postgres.NewCampaignStorage(app.DB), redis.NewDayStorage(app.Redis)),
		validator:         app.Validator,
	}
}

func (h *ModerationHandler) GetCampaignsForModeration(c fiber.Ctx) error {
	tracer := otel.Tracer("moderation-handler")
	ctx, span := tracer.Start(c.Context(), "GetCampaignsForModeration")
	defer span.End()

	var moderationDTO dto.GetCampaignsForModerationDTO

	if err := c.Bind().Query(&moderationDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(moderationDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	campaigns, err := h.ModerationService.GetCampaignsForModeration(ctx, moderationDTO)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.Int("count", len(campaigns)),
		attribute.String("endpoint", "/moderation/campaigns"),
	)

	return c.Status(fiber.StatusOK).JSON(campaigns)
}

func (h *ModerationHandler) Approve(c fiber.Ctx) error {
	tracer := otel.Tracer("moderation-handler")
	ctx, span := tracer.Start(c.Context(), "Approve")
	defer span.End()

	var approveDTO dto.ApproveCampaignDTO

	if err := c.Bind().URI(&approveDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	campaign, err := h.CampaignService.GetCampaignByIdInsecure(ctx, approveDTO.CampaignID)

	if err != nil || campaign == (dto.CampaignDTO{}) {
		span.RecordError(err)
		return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	if err := h.ModerationService.Approve(ctx, approveDTO.CampaignID); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

func (h *ModerationHandler) Reject(c fiber.Ctx) error {
	tracer := otel.Tracer("moderation-handler")
	ctx, span := tracer.Start(c.Context(), "Reject")
	defer span.End()

	var rejectDTO dto.RejectCampaignDTO

	if err := c.Bind().URI(&rejectDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	campaign, err := h.CampaignService.GetCampaignByIdInsecure(ctx, rejectDTO.CampaignID)

	if err != nil || campaign == (dto.CampaignDTO{}) {
		span.RecordError(err)
		return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	if err := h.ModerationService.Reject(ctx, rejectDTO.CampaignID); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}

func (h *ModerationHandler) Setup(router fiber.Router) {
	moderationGroup := router.Group("/moderation")

	moderationGroup.Get("/campaigns", h.GetCampaignsForModeration)
	moderationGroup.Post("/:campaignId/approve", h.Approve)
	moderationGroup.Post("/:campaignId/reject", h.Reject)
}
