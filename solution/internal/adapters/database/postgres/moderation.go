package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"solution/internal/domain/common/errorz"
)

type moderationStorage struct {
	db *pgxpool.Pool
}

// TODO turn moderation on and off

func NewModerationStorage(db *pgxpool.Pool) *moderationStorage {
	return &moderationStorage{
		db: db,
	}
}

const approve = `-- name: Approve :exec
UPDATE campaigns
SET approved = true
WHERE id = $1
`

func (s *moderationStorage) Approve(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("moderation-storage")
	ctx, span := tracer.Start(ctx, "moderation-storage")
	defer span.End()

	cmd, err := s.db.Exec(ctx, approve, id)

	if cmd.RowsAffected() == 0 {
		return errorz.NotFound
	}

	return err
}

const getCampaignsForModeration = `-- name: GetCampaignsForModeration :many
SELECT c.id,
       c.advertiser_id,
       c.ad_title,
       c.ad_text
FROM campaigns c
WHERE c.approved = false
LIMIT $1 OFFSET $2
`

type GetCampaignsForModerationParams struct {
	Limit  int32
	Offset int32
}

type GetCampaignsForModerationRow struct {
	ID           uuid.UUID
	AdvertiserID uuid.UUID
	AdTitle      string
	AdText       string
}

func (s *moderationStorage) GetCampaignsForModeration(ctx context.Context, arg GetCampaignsForModerationParams) ([]GetCampaignsForModerationRow, error) {
	tracer := otel.Tracer("moderation-storage")
	ctx, span := tracer.Start(ctx, "moderation-storage")
	defer span.End()

	rows, err := s.db.Query(ctx, getCampaignsForModeration, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCampaignsForModerationRow
	for rows.Next() {
		var i GetCampaignsForModerationRow
		if err := rows.Scan(
			&i.ID,
			&i.AdvertiserID,
			&i.AdTitle,
			&i.AdText,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reject = `-- name: Reject :exec
UPDATE campaigns
SET approved = false
WHERE id = $1
`

func (s *moderationStorage) Reject(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("moderation-storage")
	ctx, span := tracer.Start(ctx, "moderation-storage")
	defer span.End()

	cmd, err := s.db.Exec(ctx, reject, id)

	if cmd.RowsAffected() == 0 {
		return errorz.NotFound
	}

	return err
}
