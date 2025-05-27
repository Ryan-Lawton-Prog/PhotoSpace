package photo

import (
	"bytes"
	"context"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"

	"ryanlawton.art/photospace/internal/api/models"

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
	metadata.ThumbnailURL = filepath.Join(pr.UploadPath, metadata.ID+".thumbnail")

	thumbnailBlob, err := pr.generateThumbnail(blob)
	if err != nil {
		return err
	}

	pr.savePhoto(metadata.BucketURL, blob)
	pr.savePhoto(metadata.ThumbnailURL, &thumbnailBlob)

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

// FetchThumbnail fetches a photo from the database
func (pr *BucketRepository) FetchThumbnail(ctx context.Context, metadata *models.PhotoMetadata) (models.PhotoBlob, error) {
	if metadata.ThumbnailURL == "" {
		metadata.ThumbnailURL = filepath.Join(pr.UploadPath, metadata.ID+".thumbnail")
	}

	f, err := os.Open(metadata.ThumbnailURL)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		log.Printf("Attempting to generate Thumbnail", err.Error())
		blob, err2 := pr.FetchPhoto(ctx, metadata)
		if err2 != nil {
			log.Printf("Error opening file: %s", err2.Error())
			return nil, err2
		}

		thumbnailBlob, errThumnail := pr.generateThumbnail(&blob)
		if errThumnail != nil {
			return nil, errThumnail
		}

		pr.savePhoto(metadata.ThumbnailURL, &thumbnailBlob)
		log.Printf("Saved photo with path: %s", metadata.BucketURL)

		return thumbnailBlob, nil
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

func (pr *BucketRepository) savePhoto(url string, blob *models.PhotoBlob) error {
	newFile, err := os.Create(url)
	if err != nil {
		return err
	}

	defer newFile.Close() // idempotent, okay to call twice

	if _, err := newFile.Write(*blob); err != nil || newFile.Close() != nil {
		return err
	}

	return nil
}

func (pr *BucketRepository) generateThumbnail(blob *models.PhotoBlob) (models.PhotoBlob, error) {
	image, err := jpeg.Decode(bytes.NewReader(*blob))
	if err != nil {
		return nil, err
	}

	// resize image to 200x200
	thumbnail := resize.Thumbnail(200, 200, image, resize.Lanczos3)
	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, thumbnail, nil)
	thumbnailBlob := buffer.Bytes()
	return thumbnailBlob, nil
}

func (pr *BucketRepository) DeletePhoto(ctx context.Context, metadata *models.PhotoMetadata) error {
	err := os.Remove(metadata.BucketURL)
	if err != nil {
		log.Printf("Error deleting primary image: %s", err.Error())
		return err
	}

	err = os.Remove(metadata.ThumbnailURL)
	if err != nil {
		log.Printf("Error deleting thumbnail image: %s", err.Error())
		return err
	}

	return nil
}
