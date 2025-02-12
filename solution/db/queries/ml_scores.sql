-- name: InsertOrUpdateMlScore :one
INSERT INTO ml_scores (client_id, advertiser_id, score) VALUES ($1, $2, $3)
ON CONFLICT (client_id, advertiser_id)
    DO UPDATE SET score = $3
RETURNING id;