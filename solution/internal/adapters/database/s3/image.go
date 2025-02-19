package s3

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel"
	"io"
	"mime/multipart"
	"solution/internal/domain/common/errorz"
)

type imageStorage struct {
	client *minio.Client
	bucket string
}

func NewImageStorage(client *minio.Client, bucket string) *imageStorage {
	return &imageStorage{
		client: client,
		bucket: bucket,
	}
}

// TODO добавить поддержку jpg/jpeg

func (s *imageStorage) UploadImage(ctx context.Context, image *multipart.File, campaignId uuid.UUID, size int64) error {
	tracer := otel.Tracer("image-storage")
	ctx, span := tracer.Start(ctx, "image-storage")
	defer span.End()

	_, err := s.client.PutObject(ctx, s.bucket, fmt.Sprintf("%s.png", campaignId), *image, size, minio.PutObjectOptions{
		ContentType: "image/png",
	})

	return err
}

func (s *imageStorage) GetImage(ctx context.Context, campaignId uuid.UUID) ([]byte, error) {
	tracer := otel.Tracer("image-storage")
	ctx, span := tracer.Start(ctx, "image-storage")
	defer span.End()

	_, err := s.client.StatObject(ctx, s.bucket, fmt.Sprintf("%s.png", campaignId.String()), minio.StatObjectOptions{})
	if err != nil {
		return nil, errorz.NotFound
	}

	object, err := s.client.GetObject(ctx, s.bucket, fmt.Sprintf("%s.png", campaignId.String()), minio.GetObjectOptions{})
	defer func(object *minio.Object) {
		_ = object.Close()
	}(object)

	if err != nil {
		return nil, err
	}

	buffer, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (s *imageStorage) DeleteImage(ctx context.Context, campaignId uuid.UUID) error {
	tracer := otel.Tracer("image-storage")
	ctx, span := tracer.Start(ctx, "image-storage")
	defer span.End()

	return s.client.RemoveObject(ctx, s.bucket, fmt.Sprintf("%s.png", campaignId), minio.RemoveObjectOptions{})
}
