package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/redis"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type dayService interface {
	SetDay(ctx context.Context, dto dto.SetDayDTO) (int, error)
	GetDay(ctx context.Context) (int, error)
}

type DayHandler struct {
	dayService dayService
	validator  *validator.Validator
}

func NewDayHandler(app *app.App) *DayHandler {
	dayStorage := redis.NewDayStorage(app.Redis)

	return &DayHandler{
		dayService: service.NewDayService(dayStorage),
		validator:  app.Validator,
	}
}

func (h *DayHandler) SetDay(c fiber.Ctx) error {
	tracer := otel.Tracer("day-handler")
	ctx, span := tracer.Start(c.Context(), "SetDay")
	defer span.End()

	var dayDTO dto.SetDayDTO

	if err := c.Bind().Body(&dayDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(dayDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	day, err := h.dayService.SetDay(ctx, dayDTO)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SetDayDTO{CurrentDate: day})
}

func (h *DayHandler) Setup(router fiber.Router) {
	router.Post("/time/advance", h.SetDay)
}
