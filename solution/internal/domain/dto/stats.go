package dto

type GetStatsByCampaignIDDTO struct {
	CampaignID string `params:"campaignId"`
}

type GetStatsByAdvertiserIDDTO struct {
	AdvertiserID string `params:"advertiserId"`
}

type StatsDTO struct {
	ImpressionsCount int     `json:"impressions_count"`
	ClicksCount      int     `json:"clicks_count"`
	Conversion       float64 `json:"conversion"`
	SpentImpressions float64 `json:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total"`
}
