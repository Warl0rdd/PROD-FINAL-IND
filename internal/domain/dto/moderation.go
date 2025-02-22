package dto

type GetCampaignsForModerationDTO struct {
	Offset int `query:"offset" validate:"omitempty,min=0"`
	Limit  int `query:"limit" validate:"omitempty,min=1"`
}

type CampaignForModerationDTO struct {
	ID           string `json:"id" validate:"required"`
	AdvertiserID string `json:"advertiser_id" validate:"required"`
	AdTitle      string `json:"ad_title" validate:"required"`
	AdText       string `json:"ad_text" validate:"required"`
}

type ApproveCampaignDTO struct {
	CampaignID string `params:"campaignId" validate:"required,uuid"`
}

type RejectCampaignDTO struct {
	CampaignID string `params:"campaignId" validate:"required,uuid"`
}
