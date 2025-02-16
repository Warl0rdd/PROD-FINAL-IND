package dto

type LLMRequestDTO struct {
	AdvertiserName string `json:"advertiser_name" validate:"required"`
	CampaignTitle  string `json:"campaign_title" validate:"required"`
}

type LLMResponseDTO struct {
	CampaignText string `json:"campaign_text"`
}
