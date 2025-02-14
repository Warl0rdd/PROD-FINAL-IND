package v1

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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

// TODO: testing

func (h *ClientHandler) CreateClients(c fiber.Ctx) error {
	var DTOs []dto.CreateClientDTO

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

	_, err := h.clientService.CreateClient(c.Context(), DTOs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(DTOs)
}

func (h *ClientHandler) GetClientById(c fiber.Ctx) error {
	var getClientDTO dto.GetClientByIdDTO

	if err := c.Bind().URI(&getClientDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.validator.ValidateData(getClientDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	client, err := h.clientService.GetClientById(c.Context(), uuid.MustParse(getClientDTO.ClientId))
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
