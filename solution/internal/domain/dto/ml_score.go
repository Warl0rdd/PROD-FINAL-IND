package dto

type CreateMlScoreDTO struct {
	ClientID     string `json:"client_id" validate:"required"`
	AdvertiserID string `json:"advertiser_id" validate:"required"`
	Score        int    `json:"score" validate:"required"`
}
