-- name: GetStatsByCampaignID :one
SELECT c.impressions_count,
       c.clicks_count,
       COALESCE((c.clicks_count::numeric / NULLIF(c.impressions_count, 0)) * 100, 0)     AS conversion,
       c.cost_per_impression * c.impressions_count                                       AS spent_impressions,
       c.cost_per_click * c.clicks_count                                                 AS spent_clicks,
       (c.cost_per_impression * c.impressions_count + c.cost_per_click * c.clicks_count) AS spent_total
FROM campaigns c
WHERE c.id = $1;

-- name: GetStatsByAdvertiserID :one
WITH spent AS (SELECT SUM(cost_per_impression * impressions_count) AS spent_impressions,
                      SUM(cost_per_click * clicks_count)           AS spent_clicks
               FROM campaigns
               WHERE advertiser_id = $1)
SELECT SUM(impressions_count)                                                                AS total_impressions,
       SUM(clicks_count)                                                                     AS total_clicks,
       COALESCE(((SUM(clicks_count)::numeric / NULLIF(SUM(impressions_count), 0)) * 100), 0) AS conversion,
       COALESCE(SUM(spent.spent_impressions), 0)                                             AS spent_impressions,
       COALESCE(SUM(spent.spent_clicks), 0)                                                  AS spent_clicks,
       COALESCE(SUM(spent.spent_impressions) + SUM(spent.spent_clicks), 0)                   AS spent_total
FROM campaigns с
         CROSS JOIN spent
WHERE с.advertiser_id = $1;

-- name: GetDailyStatsByCampaignID :one
SELECT COALESCE(imps.impressions_count, 0)                                                 AS impressions_count,
       COALESCE(clks.clicks_count, 0)                                                      AS clicks_count,
       COALESCE((clks.clicks_count::numeric / NULLIF(imps.impressions_count, 0)) * 100, 0) AS conversion,
       (c.cost_per_impression * COALESCE(imps.impressions_count, 0))                       AS spent_impressions,
       (c.cost_per_click * COALESCE(clks.clicks_count, 0))                                 AS spent_clicks,
       (c.cost_per_impression * COALESCE(imps.impressions_count, 0)
           + c.cost_per_click * COALESCE(clks.clicks_count, 0))                            AS spent_total
FROM campaigns c
         LEFT JOIN (SELECT campaign_id, COUNT(*) AS impressions_count
                    FROM impressions impr
                    WHERE impr.day = $2
                    GROUP BY campaign_id) imps ON c.id = imps.campaign_id
         LEFT JOIN (SELECT campaign_id, COUNT(*) AS clicks_count
                    FROM clicks cl
                    WHERE cl.day = $2
                    GROUP BY campaign_id) clks ON c.id = clks.campaign_id
WHERE c.id = $1;


-- name: GetDailyStatsByAdvertiserID :one
SELECT COALESCE(SUM(imps.impressions_count), 0)                                            AS impressions_count,
       COALESCE(SUM(clks.clicks_count), 0)                                                 AS clicks_count,
       COALESCE((SUM(clks.clicks_count)::numeric / NULLIF(SUM(imps.impressions_count), 0)) * 100, 0) AS conversion,
       COALESCE(SUM(c.cost_per_impression * COALESCE(imps.impressions_count, 0)), 0)       AS spent_impressions,
       COALESCE(SUM(c.cost_per_click * COALESCE(clks.clicks_count, 0)), 0)                 AS spent_clicks,
       COALESCE(SUM(c.cost_per_impression * COALESCE(imps.impressions_count, 0))
                    + SUM(c.cost_per_click * COALESCE(clks.clicks_count, 0)), 0)           AS spent_total
FROM campaigns c
         LEFT JOIN (SELECT campaign_id, COUNT(*) AS impressions_count
                    FROM impressions impr
                    WHERE impr.day = $2
                    GROUP BY campaign_id) imps ON c.id = imps.campaign_id
         LEFT JOIN (SELECT campaign_id, COUNT(*) AS clicks_count
                    FROM clicks cl
                    WHERE cl.day = $2
                    GROUP BY campaign_id) clks ON c.id = clks.campaign_id
WHERE c.advertiser_id = $1;
