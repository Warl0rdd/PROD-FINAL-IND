package service

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"mime/multipart"
	"solution/internal/domain/dto"
)

type ImageStorage interface {
	UploadImage(ctx context.Context, image *multipart.File, campaignId uuid.UUID, size int64) error
	GetImage(ctx context.Context, campaignId uuid.UUID) ([]byte, error)
	DeleteImage(ctx context.Context, campaignId uuid.UUID) error
}

type imageService struct {
	imageStorage ImageStorage
}

func NewImageService(imageStorage ImageStorage) *imageService {
	return &imageService{
		imageStorage: imageStorage,
	}
}

func (s *imageService) UploadImage(ctx context.Context, imageDTO dto.ImageUploadDTO) error {
	tracer := otel.Tracer("image-service")
	ctx, span := tracer.Start(ctx, "image-service")
	defer span.End()

	file, err := imageDTO.Image.Open()
	if err != nil {
		span.RecordError(err)
		return err
	}

	defer file.Close()

	campaignId, err := uuid.Parse(imageDTO.CampaignID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return s.imageStorage.UploadImage(ctx, &file, campaignId, imageDTO.Image.Size)
}

func (s *imageService) GetImage(ctx context.Context, imageDTO dto.GetImageDTO) ([]byte, error) {
	tracer := otel.Tracer("image-service")
	ctx, span := tracer.Start(ctx, "image-service")
	defer span.End()

	campaignId, err := uuid.Parse(imageDTO.CampaignID)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return s.imageStorage.GetImage(ctx, campaignId)
}

func (s *imageService) DeleteImage(ctx context.Context, imageDTO dto.DeleteImageDTO) error {
	tracer := otel.Tracer("image-service")
	ctx, span := tracer.Start(ctx, "image-service")
	defer span.End()

	campaignId, err := uuid.Parse(imageDTO.CampaignID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return s.imageStorage.DeleteImage(ctx, campaignId)
}
