package v1

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/database/redis"
	"solution/internal/domain/common/errorz"
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
		adsService: service.NewAdsService(adsStorage, dayStorage, redisLearningStorage, postgresLearningStorage, postgres.NewClientStorage(app.DB)),
		validator:  app.Validator,
	}
}

func (h *AdsHandler) GetAds(c fiber.Ctx) error {
	tracer := otel.Tracer("ads-handler")
	ctx, span := tracer.Start(c.Context(), "GetAds")
	defer span.End()

	var adsDTO dto.GetAdsDTO

	if err := c.Bind().Query(&adsDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(adsDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	ads, err := h.adsService.GetAds(ctx, adsDTO)

	if err != nil {
		span.RecordError(err)
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, errorz.NotFound) {
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

	span.SetAttributes(
		attribute.String("endpoint", "/ads"),
		attribute.String("advertiserId", ads.AdvertiserID),
	)

	return c.Status(fiber.StatusOK).JSON(ads)
}

func (h *AdsHandler) Click(c fiber.Ctx) error {
	tracer := otel.Tracer("ads-handler")
	ctx, span := tracer.Start(c.Context(), "Click")
	defer span.End()

	var clickDTO dto.AddClickDTO

	if err := c.Bind().Body(&clickDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Bind().URI(&clickDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(clickDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err := h.adsService.Click(ctx, clickDTO)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, errorz.Forbidden) {
			return c.Status(fiber.StatusForbidden).JSON(dto.HTTPError{
				Code:    fiber.StatusForbidden,
				Message: "You can click on an ad only after you've seen it",
			})
		}

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
