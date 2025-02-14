package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mlScoreStorage struct {
	db *pgxpool.Pool
}

func NewMlScoreStorage(db *pgxpool.Pool) *mlScoreStorage {
	return &mlScoreStorage{
		db: db,
	}
}

const insertOrUpdateMlScore = `-- name: InsertOrUpdateMlScore :one
INSERT INTO ml_scores (client_id, advertiser_id, score) VALUES ($1, $2, $3)
ON CONFLICT (client_id, advertiser_id)
    DO UPDATE SET score = $3
RETURNING id
`

type InsertOrUpdateMlScoreParams struct {
	ClientID     uuid.UUID
	AdvertiserID uuid.UUID
	Score        float64
}

func (s *mlScoreStorage) InsertOrUpdateMlScore(ctx context.Context, arg InsertOrUpdateMlScoreParams) (uuid.UUID, error) {
	row := s.db.QueryRow(ctx, insertOrUpdateMlScore, arg.ClientID, arg.AdvertiserID, arg.Score)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
