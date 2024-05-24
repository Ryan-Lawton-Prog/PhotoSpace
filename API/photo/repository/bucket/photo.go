package photo

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"ryanlawton.art/photospace-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BucketRepository struct {
	UploadPath string
}

func NewBucketRepository(uploadPath string) *BucketRepository {
	if uploadPath == "" {
		uploadPath = os.TempDir()
	}

	return &BucketRepository{
		UploadPath: uploadPath,
	}
}

// UploadPhoto uploads a photo to the database
func (pr *BucketRepository) UploadPhoto(ctx context.Context, blob *models.PhotoBlob, metadata *models.PhotoMetadata) error {
	metadata.ID = primitive.NewObjectID().Hex()

	// Check if ID is already used
	metadata.BucketURL = filepath.Join(pr.UploadPath, metadata.ID)

	newFile, err := os.Create(metadata.BucketURL)
	if err != nil {
		return err
	}

	defer newFile.Close() // idempotent, okay to call twice

	if _, err := newFile.Write(*blob); err != nil || newFile.Close() != nil {
		return err
	}

	log.Printf("Saved photo with path: %s", metadata.BucketURL)

	return nil
}

// FetchPhoto fetches a photo from the database
func (pr *BucketRepository) FetchPhoto(ctx context.Context, metadata *models.PhotoMetadata) (models.PhotoBlob, error) {
	f, err := os.Open(metadata.BucketURL)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return nil, err
	}
	defer f.Close()

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	blob := make(models.PhotoBlob, stat.Size())
	f.Read(blob)

	return blob, nil
}
