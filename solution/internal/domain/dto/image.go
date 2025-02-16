package dto

import "mime/multipart"

type ImageUploadDTO struct {
	CampaignID string `params:"campaignId" validate:"required"`
	Image      *multipart.FileHeader
}

type GetImageDTO struct {
	CampaignID string `params:"campaignId" validate:"required"`
}

type DeleteImageDTO struct {
	CampaignID string `params:"campaignId" validate:"required"`
}

type ImageDTO struct {
	Image []byte
}
