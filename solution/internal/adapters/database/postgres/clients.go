package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"solution/internal/domain/entity"
)

type clientStorage struct {
	db *pgxpool.Pool
}

func NewClientStorage(db *pgxpool.Pool) *clientStorage {
	return &clientStorage{
		db: db,
	}
}

const createClient = `-- name: CreateClient :one
INSERT INTO clients (id, login, age, location, gender)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id)
    DO UPDATE SET login    = $2,
                  age      = $3,
                  location = $4,
                  gender   = $5
RETURNING id, login, age, location, gender
`

type CreateClientParams struct {
	ID       uuid.UUID
	Login    string
	Age      int32
	Location string
	Gender   entity.Gender
}

func (s *clientStorage) CreateClient(ctx context.Context, arg CreateClientParams) (entity.Client, error) {
	row := s.db.QueryRow(ctx, createClient,
		arg.ID,
		arg.Login,
		arg.Age,
		arg.Location,
		arg.Gender,
	)
	var i entity.Client
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Age,
		&i.Location,
		&i.Gender,
	)
	return i, err
}

const getClientById = `-- name: GetClientById :one
SELECT id, login, age, location, gender FROM clients WHERE id = $1
`

func (s *clientStorage) GetClientById(ctx context.Context, id uuid.UUID) (entity.Client, error) {
	row := s.db.QueryRow(ctx, getClientById, id)
	var i entity.Client
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Age,
		&i.Location,
		&i.Gender,
	)
	return i, err
}
