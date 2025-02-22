package entity

import "github.com/google/uuid"

type MlScore struct {
	ClientID     uuid.UUID `json:"client_id" validate:"required"`
	AdvertiserID uuid.UUID `json:"advertiser_id" validate:"required"`
	Score        int32     `json:"score" validate:"required"`
}
