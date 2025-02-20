package dto

type CampaignDTO struct {
	CampaignID        string  `json:"campaign_id" validate:"required"`
	AdvertiserID      string  `json:"advertiser_id" validate:"required"`
	ImpressionsLimit  int32   `json:"impressions_limit" validate:"required"`
	ClicksLimit       int32   `json:"clicks_limit" validate:"required"`
	CostPerImpression float64 `json:"cost_per_impression" validate:"required"`
	CostPerClick      float64 `json:"cost_per_click" validate:"required"`
	AdTitle           string  `json:"ad_title" validate:"required"`
	AdText            string  `json:"ad_text" validate:"required"`
	StartDate         int32   `json:"start_date" validate:"required"`
	EndDate           int32   `json:"end_date" validate:"required"`
	Approved          bool    `json:"approved" validate:"required"`
	Targeting         Target  `json:"targeting"`
}

type Target struct {
	Gender   string `json:"gender"`
	AgeFrom  int32  `json:"age_from" validate:"omitempty,min=0,max=200"`
	AgeTo    int32  `json:"age_to" validate:"omitempty,min=0,max=200"`
	Location string `json:"location" validate:"omitempty,max=255"`
}

type CreateCampaignDTO struct {
	AdvertiserID      string  `params:"advertiserId" validate:"required,uuid"`
	ImpressionsLimit  int32   `json:"impressions_limit" validate:"required"`
	ClicksLimit       int32   `json:"clicks_limit" validate:"required"`
	CostPerImpression float64 `json:"cost_per_impression" validate:"required"`
	CostPerClick      float64 `json:"cost_per_click" validate:"required"`
	AdTitle           string  `json:"ad_title" validate:"required"`
	AdText            string  `json:"ad_text" validate:"required"`
	StartDate         int32   `json:"start_date" validate:"required"`
	EndDate           int32   `json:"end_date" validate:"required"`
	Targeting         Target  `json:"targeting"`
}

type GetCampaignByIDDTO struct {
	CampaignID   string `params:"campaignId" validate:"required,uuid"`
	AdvertiserID string `params:"advertiserId" validate:"required,uuid"`
}

type GetCampaignsWithPaginationDTO struct {
	AdvertiserID string `params:"advertiserId" validate:"required,uuid"`
	Limit        int32  `query:"size" validate:"omitempty"`
	Page         int32  `query:"page" validate:"omitempty"`
	Offset       int32
}

type UpdateCampaignDTO struct {
	CampaignID        string       `params:"campaignId" validate:"required,uuid"`
	AdvertiserID      string       `params:"advertiserId" validate:"required,uuid"`
	ImpressionsLimit  *int32       `json:"impressions_limit"`
	ClicksLimit       *int32       `json:"clicks_limit"`
	CostPerImpression float64      `json:"cost_per_impression"`
	CostPerClick      float64      `json:"cost_per_click"`
	AdTitle           string       `json:"ad_title"`
	AdText            string       `json:"ad_text"`
	StartDate         *int         `json:"start_date"`
	EndDate           *int         `json:"end_date"`
	Targeting         TargetUpdate `json:"targeting"`
}

type TargetUpdate struct {
	Gender   *string `json:"gender"`
	AgeFrom  int32   `json:"age_from" validate:"omitempty"`
	AgeTo    int32   `json:"age_to" validate:"omitempty"`
	Location string  `json:"location" validate:"omitempty"`
}

type DeleteCampaignDTO struct {
	CampaignID   string `params:"campaignId" validate:"required,uuid"`
	AdvertiserID string `params:"advertiserId" validate:"required,uuid"`
}
