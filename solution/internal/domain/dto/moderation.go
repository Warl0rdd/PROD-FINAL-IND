package dto

type GetCampaignsForModerationDTO struct {
	Offset int `query:"offset" validate:"required"`
	Limit  int `query:"limit" validate:"required"`
}

type CampaignForModerationDTO struct {
	ID           string `json:"id" validate:"required"`
	AdvertiserID string `json:"advertiser_id" validate:"required"`
	AdTitle      string `json:"ad_title" validate:"required"`
	AdText       string `json:"ad_text" validate:"required"`
}

type ApproveCampaignDTO struct {
	CampaignID string `params:"campaignId" validate:"required"`
}

type RejectCampaignDTO struct {
	CampaignID string `params:"campaignId" validate:"required"`
}
