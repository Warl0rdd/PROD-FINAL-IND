package v1

import (
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"os"
	"solution/cmd/app"
	"solution/internal/adapters/database/s3"
	"solution/internal/domain/common/errorz"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type imageService interface {
	UploadImage(ctx context.Context, imageDTO dto.ImageUploadDTO) error
	GetImage(ctx context.Context, imageDTO dto.GetImageDTO) ([]byte, error)
	DeleteImage(ctx context.Context, imageDTO dto.DeleteImageDTO) error
}

type ImageHandler struct {
	imageService imageService
}

func NewImageHandler(app *app.App) *ImageHandler {
	return &ImageHandler{
		imageService: service.NewImageService(s3.NewImageStorage(app.Minio, os.Getenv("MINIO_IMAGES_BUCKET"))),
	}
}

func (h *ImageHandler) UploadImage(c fiber.Ctx) error {
	tracer := otel.Tracer("image-handler")
	ctx, span := tracer.Start(c.Context(), "UploadImage")
	defer span.End()

	var imageDTO dto.ImageUploadDTO

	if err := c.Bind().URI(&imageDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := uuid.Validate(imageDTO.CampaignID); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	file, err := c.FormFile("image")

	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	headerBuff := make([]byte, 512)

	content, err := file.Open()

	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	fileHeader, err := content.Read(headerBuff)
	defer content.Close()

	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	contentType := http.DetectContentType(headerBuff[:fileHeader])

	if contentType != "image/png" && contentType != "image/jpeg" {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid file type, only png and jpg/jpeg images are supported",
		})
	}

	span.SetAttributes(
		attribute.String("campaignId", imageDTO.CampaignID),
		attribute.String("filename", file.Filename),
		attribute.Int64("size", file.Size),
		attribute.String("endpoint", "/images"),
	)

	imageDTO.Image = file
	imageDTO.ContentType = contentType

	if err := h.imageService.UploadImage(ctx, imageDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *ImageHandler) GetImage(c fiber.Ctx) error {
	tracer := otel.Tracer("image-handler")
	ctx, span := tracer.Start(c.Context(), "GetImage")
	defer span.End()

	var imageDTO dto.GetImageDTO

	if err := c.Bind().URI(&imageDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := uuid.Validate(imageDTO.CampaignID); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	image, err := h.imageService.GetImage(ctx, imageDTO)

	if err != nil {
		span.RecordError(err)
		if errors.Is(err, errorz.NotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	headerBuff := make([]byte, 512)

	fileHeader, err := bytes.NewReader(image).Read(headerBuff)

	if err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	c.Set("Content-Type", http.DetectContentType(headerBuff[:fileHeader]))

	return c.Status(fiber.StatusOK).Send(image)
}

func (h *ImageHandler) DeleteImage(c fiber.Ctx) error {
	tracer := otel.Tracer("image-handler")
	ctx, span := tracer.Start(c.Context(), "DeleteImage")
	defer span.End()

	var imageDTO dto.DeleteImageDTO

	if err := c.Bind().URI(&imageDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := uuid.Validate(imageDTO.CampaignID); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := h.imageService.DeleteImage(ctx, imageDTO); err != nil {
		span.RecordError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *ImageHandler) Setup(router fiber.Router) {
	router.Post("/images/:campaignId", h.UploadImage)
	router.Get("/images/:campaignId", h.GetImage)
	router.Delete("/images/:campaignId", h.DeleteImage)
}
