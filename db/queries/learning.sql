-- name: GetImpressionsForLearning :many
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
WHERE i.used_for_learning = FALSE;

-- name: UpdateLearnedImpression :exec
UPDATE impressions
SET used_for_learning = TRUE
WHERE id = $1;
