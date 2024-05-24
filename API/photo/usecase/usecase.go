package usecase

import (
	"context"

	"ryanlawton.art/photospace-api/models"
	"ryanlawton.art/photospace-api/photo"
)

type PhotoUseCase struct {
	metadataRepo photo.MetadataRepository
	bucketRepo   photo.BucketRepository
}

func NewPhotoUseCase(metadataRepo photo.MetadataRepository, bucketRepo photo.BucketRepository) *PhotoUseCase {
	return &PhotoUseCase{
		metadataRepo: metadataRepo,
		bucketRepo:   bucketRepo,
	}
}

func (b PhotoUseCase) UploadPhoto(ctx context.Context, user *models.User, blob *models.PhotoBlob, filename string) (string, error) {
	pm := &models.PhotoMetadata{
		UserID:   user.ID,
		Filename: filename,
	}

	err := b.bucketRepo.UploadPhoto(ctx, blob, pm)
	if err != nil {
		return "Error Uploading File to Bucket", err
	}

	return b.metadataRepo.UploadPhoto(ctx, pm)
}

func (b PhotoUseCase) FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.PhotoMetadata, models.PhotoBlob, error) {
	pm := &models.PhotoMetadata{
		ID:     photoId,
		UserID: user.ID,
	}

	if err := b.metadataRepo.FetchPhoto(ctx, pm); err != nil {
		return nil, nil, err
	}
	blob, err := b.bucketRepo.FetchPhoto(ctx, pm)
	if err != nil {
		return nil, nil, err
	}

	return pm, blob, nil
}

func (b PhotoUseCase) FetchPhotoAllIDs(ctx context.Context, user *models.User) ([]string, error) {
	return b.metadataRepo.FetchPhotoAllIDs(ctx, user)
}
