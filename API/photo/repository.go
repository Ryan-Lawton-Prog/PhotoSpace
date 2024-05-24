package photo

import (
	"context"

	"ryanlawton.art/photospace-api/models"
)

type MetadataRepository interface {
	UploadPhoto(ctx context.Context, pm *models.PhotoMetadata) (string, error)
	FetchPhoto(ctx context.Context, pm *models.PhotoMetadata) error
	FetchPhotoAllIDs(ctx context.Context, user *models.User) ([]string, error)
	// FetchAlbum(ctx context.Context, user *models.User, id string) error
}

type BucketRepository interface {
	UploadPhoto(ctx context.Context, blob *models.PhotoBlob, metadata *models.PhotoMetadata) error
	FetchPhoto(ctx context.Context, metadata *models.PhotoMetadata) (models.PhotoBlob, error)
	// DeletePhoto() error
}
