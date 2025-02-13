-- name: GetEligibleAds :many
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

-- name: AddImpression :exec
INSERT INTO impressions (campaign_id, client_id, day)
VALUES ($1, $2, $3)
ON CONFLICT (campaign_id, client_id) DO NOTHING;

-- name: AddClick :exec
INSERT INTO clicks (campaign_id, client_id, day)
VALUES ($1, $2, $3)
ON CONFLICT (campaign_id, client_id) DO NOTHING;
