package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

type adsStorage struct {
	db *pgxpool.Pool
}

func NewAdsStorage(db *pgxpool.Pool) *adsStorage {
	return &adsStorage{
		db: db,
	}
}

const addClick = `-- name: AddClick :exec
INSERT INTO clicks (campaign_id, client_id, day)
VALUES ($1, $2, $3)
ON CONFLICT (campaign_id, client_id) DO NOTHING
`

type AddClickParams struct {
	CampaignID uuid.UUID
	ClientID   uuid.UUID
	Day        int32
}

func (s *adsStorage) AddClick(ctx context.Context, arg AddClickParams) error {
	tracer := otel.Tracer("ads-storage")
	ctx, span := tracer.Start(ctx, "AddClick")
	defer span.End()

	_, err := s.db.Exec(ctx, addClick, arg.CampaignID, arg.ClientID, arg.Day)
	return err
}

const addImpression = `-- name: AddImpression :exec
INSERT INTO impressions (campaign_id, client_id, day, model_score)
VALUES ($1, $2, $3, $4)
ON CONFLICT (campaign_id, client_id) DO NOTHING
`

type AddImpressionParams struct {
	CampaignID uuid.UUID
	ClientID   uuid.UUID
	Day        int32
	ModelScore float64
}

func (s *adsStorage) AddImpression(ctx context.Context, arg AddImpressionParams) error {
	tracer := otel.Tracer("ads-storage")
	ctx, span := tracer.Start(ctx, "AddImpression")
	defer span.End()

	_, err := s.db.Exec(ctx, addImpression,
		arg.CampaignID,
		arg.ClientID,
		arg.Day,
		arg.ModelScore,
	)
	return err
}

const getEligibleAds = `-- name: GetEligibleAds :many
SELECT c.id,
       c.advertiser_id,
       c.cost_per_impression,
       c.cost_per_click,
       c.ad_title,
       c.ad_text,
       ms.score
FROM campaigns c
         INNER JOIN ml_scores ms on c.advertiser_id = ms.advertiser_id AND ms.client_id = $1
         INNER JOIN clients cl ON cl.id = $1
WHERE CASE
          WHEN c.gender = 'ALL' THEN TRUE
          WHEN c.gender != 'ALL' THEN CASE
                                          WHEN c.gender = 'MALE' THEN cl.gender = 'MALE'
                                          WHEN c.gender = 'FEMALE' THEN cl.gender = 'FEMALE' END END
  AND c.age_from <= cl.age
  AND c.age_to >= cl.age
  AND CASE WHEN c.location = '' THEN TRUE WHEN c.location != 'ALL' THEN cl.location = c.location END
  AND c.start_date <= $2
  AND c.end_date >= $2
  AND c.clicks_count < c.clicks_limit
  AND c.impressions_count < c.impression_limit
  AND NOT EXISTS (SELECT 1
                  FROM impressions i
                  WHERE i.campaign_id = c.id
                    AND i.client_id = $1);
`

type GetEligibleAdsParams struct {
	ClientID uuid.UUID
	Day      int32
}

type GetEligibleAdsRow struct {
	ID                uuid.UUID
	AdvertiserID      uuid.UUID
	CostPerImpression float64
	CostPerClick      float64
	AdTitle           string
	AdText            string
	Score             float64
}

func (s *adsStorage) GetEligibleAds(ctx context.Context, arg GetEligibleAdsParams) ([]GetEligibleAdsRow, error) {
	tracer := otel.Tracer("ads-storage")
	ctx, span := tracer.Start(ctx, "GetEligibleAds")
	defer span.End()

	rows, err := s.db.Query(ctx, getEligibleAds, arg.ClientID, arg.Day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEligibleAdsRow
	for rows.Next() {
		var i GetEligibleAdsRow
		if err := rows.Scan(
			&i.ID,
			&i.AdvertiserID,
			&i.CostPerImpression,
			&i.CostPerClick,
			&i.AdTitle,
			&i.AdText,
			&i.Score,
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
