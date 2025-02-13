package dto

type GetAdsDTO struct {
	ClientID string `query:"client_id" validate:"required"`
}

type AdDTO struct {
	AdID         string `json:"ad_id"`
	AdvertiserID string `json:"advertiser_id"`
	AdTitle      string `json:"ad_title"`
	AdText       string `json:"ad_text"`
}
