package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
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
	tracer := otel.Tracer("advertiser-handler")
	ctx, span := tracer.Start(c.Context(), "CreateAdvertiser")
	defer span.End()

	var DTOs []dto.CreateAdvertiserDTO

	if err := c.Bind().Body(&DTOs); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	for _, v := range DTOs {
		if err := h.validator.ValidateData(v); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	span.SetAttributes(
		attribute.Int("advertisers.count", len(DTOs)),
		attribute.String("endpoint", "/advertisers/bulk"),
	)

	_, err := h.advertiserService.CreateAdvertiser(ctx, DTOs[0])
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(DTOs)
}

func (h *AdvertiserHandler) GetAdvertiserById(c fiber.Ctx) error {
	tracer := otel.Tracer("advertiser-handler")
	ctx, span := tracer.Start(c.Context(), "GetAdvertiserById")
	defer span.End()

	var getAdvertiserDTO dto.GetAdvertiserByIdDTO

	if err := c.Bind().URI(&getAdvertiserDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(getAdvertiserDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.String("advertiserId", getAdvertiserDTO.AdvertiserID),
		attribute.String("endpoint", "/advertisers/{advertiserId}"),
	)

	advertiser, err := h.advertiserService.GetAdvertiserById(ctx, getAdvertiserDTO)
	if err != nil {
		span.RecordError(err)
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
