package usecase

import (
	"context"

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

func (b PhotoUseCase) UploadPhoto(ctx context.Context, user *models.User, photo []byte, filename string) error {
	pm := &models.Photo{
		Photo:    photo,
		Filename: filename,
	}

	return b.photoRepo.UploadPhoto(ctx, user, pm)
}

func (b PhotoUseCase) FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.Photo, error) {
	pm := &models.Photo{
		ID: photoId,
	}
	return b.photoRepo.FetchPhoto(ctx, user, pm)
}
