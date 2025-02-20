package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"solution/internal/domain/entity"
)

type campaignStorage struct {
	db *pgxpool.Pool
}

func NewCampaignStorage(db *pgxpool.Pool) *campaignStorage {
	return &campaignStorage{
		db: db,
	}
}

const createCampaign = `-- name: CreateCampaign :one
INSERT INTO campaigns (advertiser_id, impression_limit, clicks_limit, cost_per_impression, cost_per_click, ad_title,
                       ad_text, start_date, end_date, age_from, age_to, location, gender)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13::campaign_gender)
RETURNING id, advertiser_id;
`

type CreateCampaignParams struct {
	AdvertiserID      uuid.UUID
	ImpressionLimit   int32
	ClicksLimit       int32
	CostPerImpression float64
	CostPerClick      float64
	AdTitle           string
	AdText            string
	StartDate         int32
	EndDate           int32
	AgeFrom           pgtype.Int4
	AgeTo             pgtype.Int4
	Location          pgtype.Text
	Gender            entity.CampaignGender
}

type CreateCampaignRow struct {
	ID           uuid.UUID
	AdvertiserID uuid.UUID
}

func (s *campaignStorage) CreateCampaign(ctx context.Context, arg CreateCampaignParams) (CreateCampaignRow, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	row := s.db.QueryRow(ctx, createCampaign,
		arg.AdvertiserID,
		arg.ImpressionLimit,
		arg.ClicksLimit,
		arg.CostPerImpression,
		arg.CostPerClick,
		arg.AdTitle,
		arg.AdText,
		arg.StartDate,
		arg.EndDate,
		arg.AgeFrom,
		arg.AgeTo,
		arg.Location,
		arg.Gender,
	)
	var i CreateCampaignRow
	err := row.Scan(&i.ID, &i.AdvertiserID)
	return i, err
}

const deleteCampaign = `-- name: DeleteCampaign :exec
DELETE
FROM campaigns
WHERE id = $1
  AND advertiser_id = $2
`

type DeleteCampaignParams struct {
	ID           uuid.UUID
	AdvertiserID uuid.UUID
}

func (s *campaignStorage) DeleteCampaign(ctx context.Context, arg DeleteCampaignParams) error {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	_, err := s.db.Exec(ctx, deleteCampaign, arg.ID, arg.AdvertiserID)
	return err
}

const getCampaignById = `-- name: GetCampaignById :one
SELECT c.id,
       c.advertiser_id,
       c.impression_limit,
       c.clicks_limit,
       c.cost_per_impression,
       c.cost_per_click,
       c.ad_title,
       c.ad_text,
       c.start_date,
       c.end_date,
       c.gender,
       c.age_from,
       c.age_to,
       c.location,
	   c.approved
FROM campaigns c
WHERE c.advertiser_id = $1
  AND c.id = $2
`

type GetCampaignByIdParams struct {
	AdvertiserID uuid.UUID
	ID           uuid.UUID
}

func (s *campaignStorage) GetCampaignById(ctx context.Context, arg GetCampaignByIdParams) (entity.Campaign, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	row := s.db.QueryRow(ctx, getCampaignById, arg.AdvertiserID, arg.ID)
	var i entity.Campaign
	err := row.Scan(
		&i.ID,
		&i.AdvertiserID,
		&i.ImpressionLimit,
		&i.ClicksLimit,
		&i.CostPerImpression,
		&i.CostPerClick,
		&i.AdTitle,
		&i.AdText,
		&i.StartDate,
		&i.EndDate,
		&i.Gender,
		&i.AgeFrom,
		&i.AgeTo,
		&i.Location,
		&i.Approved,
	)
	return i, err
}

const getCampaignByIdInsecure = `-- name: GetCampaignById :one
SELECT c.id,
       c.advertiser_id,
       c.impression_limit,
       c.clicks_limit,
       c.cost_per_impression,
       c.cost_per_click,
       c.ad_title,
       c.ad_text,
       c.start_date,
       c.end_date,
       c.gender,
       c.age_from,
       c.age_to,
       c.location,
	   c.approved
FROM campaigns c
WHERE c.id = $1
`

func (s *campaignStorage) GetCampaignByIdInsecure(ctx context.Context, campaignId uuid.UUID) (entity.Campaign, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	row := s.db.QueryRow(ctx, getCampaignByIdInsecure, campaignId)
	var i entity.Campaign
	err := row.Scan(
		&i.ID,
		&i.AdvertiserID,
		&i.ImpressionLimit,
		&i.ClicksLimit,
		&i.CostPerImpression,
		&i.CostPerClick,
		&i.AdTitle,
		&i.AdText,
		&i.StartDate,
		&i.EndDate,
		&i.Gender,
		&i.AgeFrom,
		&i.AgeTo,
		&i.Location,
		&i.Approved,
	)
	return i, err
}

const getCampaignWithPagination = `-- name: GetCampaignWithPagination :many
SELECT c.id,
       c.advertiser_id,
       c.impression_limit,
       c.clicks_limit,
       c.cost_per_impression,
       c.cost_per_click,
       c.ad_title,
       c.ad_text,
       c.start_date,
       c.end_date,
       c.gender,
       c.age_from,
       c.age_to,
       c.location,
	   c.approved
FROM campaigns c
WHERE c.advertiser_id = $1
LIMIT $2 OFFSET $3
`

type GetCampaignWithPaginationParams struct {
	AdvertiserID uuid.UUID
	Limit        int32
	Offset       int32
}

func (s *campaignStorage) GetCampaignWithPagination(ctx context.Context, arg GetCampaignWithPaginationParams) ([]entity.Campaign, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	rows, err := s.db.Query(ctx, getCampaignWithPagination, arg.AdvertiserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Campaign
	for rows.Next() {
		var i entity.Campaign
		if err := rows.Scan(
			&i.ID,
			&i.AdvertiserID,
			&i.ImpressionLimit,
			&i.ClicksLimit,
			&i.CostPerImpression,
			&i.CostPerClick,
			&i.AdTitle,
			&i.AdText,
			&i.StartDate,
			&i.EndDate,
			&i.Gender,
			&i.AgeFrom,
			&i.AgeTo,
			&i.Location,
			&i.Approved,
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

const updateCampaign = `-- name: UpdateCampaign :one
UPDATE campaigns
SET cost_per_impression = CASE WHEN $3::float != 0 THEN $3 ELSE cost_per_impression END,
    cost_per_click      = CASE WHEN $4::float != 0 THEN $4 ELSE cost_per_click END,
    ad_title            = CASE WHEN $5::text != '' THEN $5 ELSE ad_title END,
    ad_text             = CASE WHEN $6::text != '' THEN $6 ELSE ad_text END,
    gender              = CASE WHEN $7::campaign_gender != 'ALL' THEN $7::campaign_gender WHEN $7 IS NULL THEN gender ELSE 'ALL' END,
    age_from            = CASE WHEN $8::int != 0 THEN $8 ELSE age_from END,
    age_to              = CASE WHEN $9::int != 0 THEN $9 ELSE age_to END,
    location            = CASE WHEN $10::text != '' THEN $10 ELSE location END,
    impression_limit    = COALESCE($11, impression_limit),
    clicks_limit        = COALESCE($12, clicks_limit),
    start_date          = COALESCE($13, start_date),
    end_date            = COALESCE($14, end_date)
WHERE id = $1
  AND advertiser_id = $2
RETURNING id, advertiser_id, impression_limit, clicks_limit, cost_per_impression, cost_per_click, ad_title, ad_text, start_date, end_date, gender, age_from, age_to, location, approved;
`

type UpdateCampaignParams struct {
	ID                uuid.UUID
	AdvertiserID      uuid.UUID
	CostPerImpression float64
	CostPerClick      float64
	AdTitle           string
	AdText            string
	Gender            *string
	AgeFrom           pgtype.Int4
	AgeTo             pgtype.Int4
	Location          pgtype.Text
	ImpressionLimit   *int32
	ClicksLimit       *int32
	StartDate         *int
	EndDate           *int
}

func (s *campaignStorage) UpdateCampaign(ctx context.Context, arg UpdateCampaignParams) (entity.Campaign, error) {
	tracer := otel.Tracer("campaign-service")
	ctx, span := tracer.Start(ctx, "campaign-service")
	defer span.End()

	row := s.db.QueryRow(ctx, updateCampaign,
		arg.ID,
		arg.AdvertiserID,
		arg.CostPerImpression,
		arg.CostPerClick,
		arg.AdTitle,
		arg.AdText,
		arg.Gender,
		arg.AgeFrom,
		arg.AgeTo,
		arg.Location,
		arg.ImpressionLimit,
		arg.ClicksLimit,
		arg.StartDate,
		arg.EndDate,
	)
	var i entity.Campaign
	err := row.Scan(
		&i.ID,
		&i.AdvertiserID,
		&i.ImpressionLimit,
		&i.ClicksLimit,
		&i.CostPerImpression,
		&i.CostPerClick,
		&i.AdTitle,
		&i.AdText,
		&i.StartDate,
		&i.EndDate,
		&i.Gender,
		&i.AgeFrom,
		&i.AgeTo,
		&i.Location,
		&i.Approved,
	)
	return i, err
}
