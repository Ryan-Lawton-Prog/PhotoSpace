package photo

import (
	"context"
	"image"

	"ryanlawton.art/photospace-api/models"
)

type UseCase interface {
	UploadPhoto(ctx context.Context, user *models.User, photo image.Image, extension string) error
	FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.Photo, error)
	// FetchAlbum(ctx context.Context, user *models.User, id string) error
}
