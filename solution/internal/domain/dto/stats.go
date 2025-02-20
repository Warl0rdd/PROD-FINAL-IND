package dto

type GetStatsByCampaignIDDTO struct {
	CampaignID string `params:"campaignId" validate:"required,uuid"`
}

type GetStatsByAdvertiserIDDTO struct {
	AdvertiserID string `params:"advertiserId" validate:"required,uuid"`
}

type StatsDTO struct {
	ImpressionsCount int     `json:"impressions_count"`
	ClicksCount      int     `json:"clicks_count"`
	Conversion       float64 `json:"conversion"`
	SpentImpressions float64 `json:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total"`
	Day              int     `json:"day,omitempty"`
}
