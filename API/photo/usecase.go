package photo

import (
	"context"

	"ryanlawton.art/photospace-api/models"
)

type UseCase interface {
	UploadPhoto(ctx context.Context, user *models.User, photo []byte, filename string) error
	FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.Photo, error)
	// FetchAlbum(ctx context.Context, user *models.User, id string) error
}
