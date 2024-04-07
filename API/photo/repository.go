package photo

import (
	"context"

	"ryanlawton.art/photospace-api/models"
)

type Repository interface {
	UploadPhoto(ctx context.Context, user *models.User, pm *models.Photo) error
	FetchPhoto(ctx context.Context, user *models.User, pm *models.Photo) (*models.Photo, error)
	// FetchAlbum(ctx context.Context, user *models.User, id string) error
}
