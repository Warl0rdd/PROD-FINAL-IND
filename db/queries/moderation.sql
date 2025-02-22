-- name: GetCampaignsForModeration :many
SELECT c.id,
       c.advertiser_id,
       c.ad_title,
       c.ad_text
FROM campaigns c
WHERE c.approved = false
LIMIT $1 OFFSET $2;

-- name: Approve :exec
UPDATE campaigns
SET approved = true
WHERE id = $1;

-- name: Reject :exec
UPDATE campaigns
SET approved = false
WHERE id = $1;