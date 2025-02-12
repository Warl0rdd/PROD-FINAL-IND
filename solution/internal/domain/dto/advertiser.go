package dto

type CreateAdvertiserDTO struct {
	AdvertiserID string `json:"advertiser_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
}

type GetAdvertiserByIdDTO struct {
	AdvertiserID string `params:"advertiserId" validate:"required"`
}
