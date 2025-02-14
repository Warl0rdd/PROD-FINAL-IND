package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type learningStorage struct {
	db *pgxpool.Pool
}

func NewLearningStorage(db *pgxpool.Pool) *learningStorage {
	return &learningStorage{
		db: db,
	}
}

const getImpressionsForLearning = `-- name: GetImpressionsForLearning :many
SELECT i.id,
       i.used_for_learning,
       i.model_score,
       CASE
           WHEN c.id IS NOT NULL THEN TRUE
           ELSE FALSE
           END  AS clicked_after,
       ms.score AS score
FROM impressions i
         LEFT JOIN clicks c
                   ON c.campaign_id = i.campaign_id
                       AND c.client_id = i.client_id
         INNER JOIN public.campaigns cmp ON i.campaign_id = cmp.id
         INNER JOIN ml_scores ms ON ms.client_id = i.client_id AND ms.advertiser_id = cmp.advertiser_id
WHERE i.used_for_learning = FALSE
`

type GetImpressionsForLearningRow struct {
	ID              uuid.UUID
	UsedForLearning bool
	ModelScore      float64
	ClickedAfter    bool
	Score           float64
}

func (s *learningStorage) GetImpressionsForLearning(ctx context.Context) ([]GetImpressionsForLearningRow, error) {
	rows, err := s.db.Query(ctx, getImpressionsForLearning)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetImpressionsForLearningRow
	for rows.Next() {
		var i GetImpressionsForLearningRow
		if err := rows.Scan(
			&i.ID,
			&i.UsedForLearning,
			&i.ModelScore,
			&i.ClickedAfter,
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

const updateLearnedImpression = `-- name: UpdateLearnedImpression :exec
UPDATE impressions
SET used_for_learning = TRUE
WHERE id = $1
`

func (s *learningStorage) UpdateLearnedImpression(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, updateLearnedImpression, id)
	return err
}
