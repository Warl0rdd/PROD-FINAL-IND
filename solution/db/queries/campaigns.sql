-- name: CreateCampaign :one
INSERT INTO campaigns (advertiser_id, impression_limit, clicks_limit, cost_per_impression, cost_per_click, ad_title,
                       ad_text, start_date, end_date, age_from, age_to, location, gender)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13::campaign_gender)
RETURNING id, advertiser_id;

-- name: GetCampaignWithPagination :many
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
       c.location
FROM campaigns c
WHERE c.advertiser_id = $1
LIMIT $2 OFFSET $3;

-- name: GetCampaignById :one
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
       c.location
FROM campaigns c
WHERE c.advertiser_id = $1
  AND c.id = $2;

-- name: UpdateCampaign :one
UPDATE campaigns
SET cost_per_impression = CASE WHEN $3::float != 0 THEN $3 ELSE cost_per_impression END,
    cost_per_click      = CASE WHEN $4::float != 0 THEN $4 ELSE cost_per_click END,
    ad_title            = CASE WHEN $5::text != '' THEN $5 ELSE ad_title END,
    ad_text             = CASE WHEN $6::text != '' THEN $6 ELSE ad_text END,
    gender              = CASE WHEN $7::campaign_gender != 'ALL' THEN $7 ELSE gender END,
    age_from            = CASE WHEN $8::int != 0 THEN $8 ELSE age_from END,
    age_to              = CASE WHEN $9::int != 0 THEN $9 ELSE age_to END,
    location            = CASE WHEN $10::text != '' THEN $10 ELSE location END,
    impression_limit    = CASE WHEN $11::int != 0 THEN $11 ELSE impression_limit END,
    clicks_limit        = CASE WHEN $12::int != 0 THEN $12 ELSE clicks_limit END
WHERE id = $1
  AND advertiser_id = $2
RETURNING *;

-- name: DeleteCampaign :exec
DELETE
FROM campaigns
WHERE id = $1
  AND advertiser_id = $2;

