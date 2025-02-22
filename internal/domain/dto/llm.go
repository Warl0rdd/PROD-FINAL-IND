package dto

type LLMRequestDTO struct {
	AdvertiserName string `json:"advertiser_name" validate:"required,max=500"`
	CampaignTitle  string `json:"campaign_title" validate:"required,max=500"`
}

type LLMResponseDTO struct {
	CampaignText string `json:"campaign_text"`
}
