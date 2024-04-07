package usecase

import (
	"context"
	"image"

	"ryanlawton.art/photospace-api/models"
	"ryanlawton.art/photospace-api/photo"
)

type PhotoUseCase struct {
	photoRepo photo.Repository
}

func NewPhotoUseCase(photoRepo photo.Repository) *PhotoUseCase {
	return &PhotoUseCase{
		photoRepo: photoRepo,
	}
}

func (b PhotoUseCase) UploadPhoto(ctx context.Context, user *models.User, photo image.Image, extension string) error {
	pm := &models.Photo{
		Photo:     photo,
		Extension: extension,
	}

	return b.photoRepo.UploadPhoto(ctx, user, pm)
}

func (b PhotoUseCase) FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.Photo, error) {
	pm := &models.Photo{
		ID: photoId,
	}
	return b.photoRepo.FetchPhoto(ctx, user, pm)
}
