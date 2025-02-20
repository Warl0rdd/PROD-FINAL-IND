package dto

import "mime/multipart"

type ImageUploadDTO struct {
	CampaignID  string `params:"campaignId" validate:"required,uuid"`
	Image       *multipart.FileHeader
	ContentType string
}

type GetImageDTO struct {
	CampaignID string `params:"campaignId" validate:"required,uuid"`
}

type DeleteImageDTO struct {
	CampaignID string `params:"campaignId" validate:"required,uuid"`
}

type ImageDTO struct {
	Image []byte
}
