-- name: GetStatsByCampaignID :one
WITH impr_stats AS (SELECT campaign_id,
                           SUM(cost) AS spent_impressions
                    FROM impressions
                    WHERE campaign_id = $1
                    GROUP BY campaign_id),
     click_stats AS (SELECT campaign_id,
                            SUM(cost) AS spent_clicks
                     FROM clicks
                     WHERE campaign_id = $1
                     GROUP BY campaign_id)
SELECT c.impressions_count,
       c.clicks_count,
       COALESCE((c.clicks_count::numeric / NULLIF(c.impressions_count, 0)) * 100, 0)     AS conversion,
       COALESCE(impr_stats.spent_impressions, 0)                                         AS spent_impressions,
       COALESCE(click_stats.spent_clicks, 0)                                             AS spent_clicks,
       COALESCE(impr_stats.spent_impressions, 0) + COALESCE(click_stats.spent_clicks, 0) AS spent_total
FROM campaigns c
         LEFT JOIN impr_stats ON c.id = impr_stats.campaign_id
         LEFT JOIN click_stats ON c.id = click_stats.campaign_id
WHERE c.id = $1;

-- name: GetStatsByAdvertiserID :one
WITH impr_stats AS (SELECT c.advertiser_id,
                           COUNT(*)    AS total_impressions,
                           SUM(i.cost) AS spent_impressions
                    FROM impressions i
                             JOIN campaigns c ON i.campaign_id = c.id
                    WHERE c.advertiser_id = $1
                    GROUP BY c.advertiser_id),
     click_stats AS (SELECT c.advertiser_id,
                            COUNT(*)     AS total_clicks,
                            SUM(c2.cost) AS spent_clicks
                     FROM clicks c2
                              JOIN campaigns c ON c2.campaign_id = c.id
                     WHERE c.advertiser_id = $1
                     GROUP BY c.advertiser_id)
SELECT COALESCE(impr_stats.total_impressions, 0)                                         AS total_impressions,
       COALESCE(click_stats.total_clicks, 0)                                             AS total_clicks,
       COALESCE(
               (click_stats.total_clicks::numeric / NULLIF(impr_stats.total_impressions, 0)) * 100,
               0
       )                                                                                 AS conversion,
       COALESCE(impr_stats.spent_impressions, 0)                                         AS spent_impressions,
       COALESCE(click_stats.spent_clicks, 0)                                             AS spent_clicks,
       COALESCE(impr_stats.spent_impressions, 0) + COALESCE(click_stats.spent_clicks, 0) AS spent_total
FROM impr_stats
         FULL OUTER JOIN click_stats USING (advertiser_id);

-- name: GetDailyStatsByCampaignID :many
SELECT COALESCE(imps.impressions_count, 0)                                  AS impressions_count,
       COALESCE(clks.clicks_count, 0)                                       AS clicks_count,
       COALESCE(
               (clks.clicks_count::numeric / NULLIF(imps.impressions_count, 0)) * 100,
               0
       )                                                                    AS conversion,
       COALESCE(imps.spent_impressions, 0)                                  AS spent_impressions,
       COALESCE(clks.spent_clicks, 0)                                       AS spent_clicks,
       COALESCE(imps.spent_impressions, 0) + COALESCE(clks.spent_clicks, 0) AS spent_total,
       COALESCE(imps.day, clks.day)                                         AS day
FROM campaigns c
         LEFT JOIN (SELECT campaign_id,
                           day,
                           COUNT(*)  AS impressions_count,
                           SUM(cost) AS spent_impressions
                    FROM impressions
                    WHERE campaign_id = $1
                    GROUP BY campaign_id, day) imps ON c.id = imps.campaign_id
         LEFT JOIN (SELECT campaign_id,
                           day,
                           COUNT(*)  AS clicks_count,
                           SUM(cost) AS spent_clicks
                    FROM clicks
                    WHERE campaign_id = $1
                    GROUP BY campaign_id, day) clks ON c.id = clks.campaign_id
WHERE c.id = $1
  AND COALESCE(imps.day, clks.day) IS NOT NULL
ORDER BY COALESCE(imps.day, clks.day);


-- name: GetDailyStatsByAdvertiserID :many
SELECT COALESCE(SUM(imps.impressions_count), 0)                          AS impressions_count,
       COALESCE(SUM(clks.clicks_count), 0)                               AS clicks_count,
       COALESCE(
               (SUM(clks.clicks_count)::numeric / NULLIF(SUM(imps.impressions_count), 0)) * 100,
               0
       )                                                                 AS conversion,
       COALESCE(SUM(imps.spent_impressions), 0)                          AS spent_impressions,
       COALESCE(SUM(clks.spent_clicks), 0)                               AS spent_clicks,
       COALESCE(SUM(imps.spent_impressions) + SUM(clks.spent_clicks), 0) AS spent_total,
       COALESCE(imps.day, clks.day)                                      AS day
FROM campaigns c
         LEFT JOIN (SELECT campaign_id,
                           day,
                           COUNT(*)  AS impressions_count,
                           SUM(cost) AS spent_impressions
                    FROM impressions
                    GROUP BY campaign_id, day) imps ON c.id = imps.campaign_id
         LEFT JOIN (SELECT campaign_id,
                           day,
                           COUNT(*)  AS clicks_count,
                           SUM(cost) AS spent_clicks
                    FROM clicks
                    GROUP BY campaign_id, day) clks ON c.id = clks.campaign_id
WHERE c.advertiser_id = $1
  AND COALESCE(imps.day, clks.day) IS NOT NULL
GROUP BY COALESCE(imps.day, clks.day)
ORDER BY COALESCE(imps.day, clks.day);
