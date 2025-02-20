package dto

type CreateMlScoreDTO struct {
	ClientID     string `json:"client_id" validate:"required,uuid"`
	AdvertiserID string `json:"advertiser_id" validate:"required,uuid"`
	Score        int    `json:"score" validate:"required,min=0"`
}
