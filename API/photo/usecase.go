package photo

import (
	"context"

	"ryanlawton.art/photospace-api/models"
)

type UseCase interface {
	UploadPhoto(ctx context.Context, user *models.User, blob *models.PhotoBlob, filename string) (string, error)
	FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.PhotoMetadata, models.PhotoBlob, error)
	FetchPhotoAllIDs(ctx context.Context, user *models.User) ([]string, error)
	// FetchAlbum(ctx context.Context, user *models.User, id string) error
}
