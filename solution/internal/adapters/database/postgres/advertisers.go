package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"solution/internal/domain/entity"
)

type advertiserStorage struct {
	db *pgxpool.Pool
}

func NewAdvertiserStorage(db *pgxpool.Pool) *advertiserStorage {
	return &advertiserStorage{
		db: db,
	}
}

const createAdvertiser = `-- name: CreateAdvertiser :one
INSERT INTO advertisers (id, name)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET name = $2
RETURNING id, name
`

type CreateAdvertiserParams struct {
	ID   uuid.UUID
	Name string
}

func (s *advertiserStorage) CreateAdvertiser(ctx context.Context, arg CreateAdvertiserParams) (entity.Advertiser, error) {
	tracer := otel.Tracer("advertiser-storage")
	ctx, span := tracer.Start(ctx, "advertiser-storage")
	defer span.End()

	row := s.db.QueryRow(ctx, createAdvertiser, arg.ID, arg.Name)
	var i entity.Advertiser
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getAdvertiserById = `-- name: GetAdvertiserById :one
SELECT id, name
FROM advertisers
WHERE id = $1
`

func (s *advertiserStorage) GetAdvertiserById(ctx context.Context, id uuid.UUID) (entity.Advertiser, error) {
	tracer := otel.Tracer("advertiser-storage")
	ctx, span := tracer.Start(ctx, "advertiser-storage")
	defer span.End()

	row := s.db.QueryRow(ctx, getAdvertiserById, id)
	var i entity.Advertiser
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}
