package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"solution/internal/adapters/database/postgres"
	"solution/internal/domain/dto"
	"solution/internal/domain/entity"
)

type ClientStorage interface {
	CreateClient(ctx context.Context, arg postgres.CreateClientParams) (entity.Client, error)
	GetClientById(ctx context.Context, id uuid.UUID) (entity.Client, error)
}

type clientService struct {
	clientStorage ClientStorage
}

func NewClientService(clientStorage ClientStorage) *clientService {
	return &clientService{
		clientStorage: clientStorage,
	}
}

func (s *clientService) CreateClient(ctx context.Context, dto []dto.CreateClientDTO) ([]entity.Client, error) {
	result := make([]entity.Client, len(dto))

	tracer := otel.Tracer("client-service")
	ctx, span := tracer.Start(ctx, "client-service")
	defer span.End()

	span.SetAttributes(attribute.Int("clients.count", len(dto)))

	for _, d := range dto {
		client, err := s.clientStorage.CreateClient(ctx, postgres.CreateClientParams{
			ID:       uuid.MustParse(d.ClientId),
			Login:    d.Login,
			Age:      int32(d.Age),
			Location: d.Location,
			Gender:   entity.Gender(d.Gender),
		})
		if err != nil {
			return nil, err
		}
		result = append(result, client)
	}

	return result, nil
}

func (s *clientService) GetClientById(ctx context.Context, id uuid.UUID) (entity.Client, error) {
	tracer := otel.Tracer("client-service")
	ctx, span := tracer.Start(ctx, "client-service")
	defer span.End()

	span.SetAttributes(attribute.String("client_id", id.String()))

	return s.clientStorage.GetClientById(ctx, id)
}
