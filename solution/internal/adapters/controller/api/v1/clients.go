package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/cmd/app"
	"solution/internal/adapters/controller/api/validator"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/logger"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
	"solution/internal/domain/service"
)

type ClientService interface {
	CreateClient(ctx context.Context, dto []dto.CreateClientDTO) ([]entity.Client, error)
	GetClientById(ctx context.Context, id uuid.UUID) (entity.Client, error)
}
type ClientHandler struct {
	clientService ClientService
	validator     *validator.Validator
}

func NewClientHandler(app *app.App) *ClientHandler {
	clientStorage := postgres.NewClientStorage(app.DB)

	return &ClientHandler{
		clientService: service.NewClientService(clientStorage),
		validator:     app.Validator,
	}
}

func (h *ClientHandler) CreateClients(c fiber.Ctx) error {
	tracer := otel.Tracer("client-handler")
	ctx, span := tracer.Start(c.Context(), "CreateClients")
	defer span.End()

	var DTOs []dto.CreateClientDTO

	if err := c.Bind().Body(&DTOs); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.Int("clients.count", len(DTOs)),
		attribute.String("endpoint", "/clients/bulk"),
	)

	logger.Log.Debugf("DTOs: %v", DTOs)

	for _, v := range DTOs {
		if err := h.validator.ValidateData(v); err != nil {
			span.RecordError(err)
			return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}

	_, err := h.clientService.CreateClient(ctx, DTOs)
	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(DTOs)
}

func (h *ClientHandler) GetClientById(c fiber.Ctx) error {
	tracer := otel.Tracer("client-handler")
	ctx, span := tracer.Start(c.Context(), "GetClientById")
	defer span.End()

	var getClientDTO dto.GetClientByIdDTO

	if err := c.Bind().URI(&getClientDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(getClientDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	span.SetAttributes(
		attribute.String("clientId", getClientDTO.ClientId),
		attribute.String("endpoint", "/clients/{clientId}"),
	)

	client, err := h.clientService.GetClientById(ctx, uuid.MustParse(getClientDTO.ClientId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.CreateClientDTO{
		ClientId: client.ID.String(),
		Login:    client.Login,
		Age:      int(client.Age),
		Location: client.Location,
		Gender:   string(client.Gender),
	})
}

func (h *ClientHandler) Setup(router fiber.Router) {
	clientGroup := router.Group("/clients")

	clientGroup.Post("/bulk", h.CreateClients)
	clientGroup.Get("/:clientId", h.GetClientById)
}
