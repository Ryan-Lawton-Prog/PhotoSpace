package photo

import (
	"context"
	"log"
	"os"

	"github.com/spf13/viper"
	"ryanlawton.art/photospace-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var uploadPath string

type PhotoMetadata struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	AlbumID   primitive.ObjectID `bson:"album_id"`
	Filename  string             `bson:"filename"`
	BucketURL string             `bson:"bucket_url"`
}

type PhotoRepository struct {
	db *mongo.Collection
}

func NewPhotoRepository(db *mongo.Database, collection string) *PhotoRepository {
	uploadPath = viper.GetString("photo.upload_dir")

	if uploadPath == "" {
		uploadPath = os.TempDir()
	}

	return &PhotoRepository{
		db: db.Collection(collection),
	}
}

// UploadPhoto uploads a photo to the database
func (pr *PhotoRepository) UploadPhoto(ctx context.Context, pm *models.PhotoMetadata) (string, error) {
	model := toModel(pm)

	res, err := pr.db.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	log.Printf("Inserted photo with ID %s", res.InsertedID.(primitive.ObjectID).Hex())

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// FetchPhoto fetches a photo from the database
func (pr *PhotoRepository) FetchPhoto(ctx context.Context, pm *models.PhotoMetadata) error {
	model := toModel(pm)

	res, err := pr.db.Find(ctx, bson.D{{Key: "user_id", Value: model.UserID}, {Key: "_id", Value: model.ID}})
	if err != nil {
		return err
	}

	defer res.Close(ctx)

	p := new(PhotoMetadata)

	for res.Next(ctx) {
		err := res.Decode(p)
		if err != nil {
			log.Printf("Error decoding photo metadata: %s", err.Error())
			return err
		}
	}

	*pm = *toPhoto(p)

	return nil
}

// FetchPhotoAllIDs fetches all photo IDs from the database
func (pr *PhotoRepository) FetchPhotoAllIDs(ctx context.Context, user *models.User) ([]string, error) {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	res, err := pr.db.Find(ctx, bson.D{{"user_id", uid}})
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	var photoIDs []string

	for res.Next(ctx) {
		photo := new(PhotoMetadata)
		err := res.Decode(photo)
		if err != nil {
			log.Printf("Error decoding photo: %s", err.Error())
			return nil, err
		}

		photoIDs = append(photoIDs, photo.ID.Hex())
	}

	return photoIDs, nil
}

func toModel(pm *models.PhotoMetadata) *PhotoMetadata {
	id, _ := primitive.ObjectIDFromHex(pm.ID)
	uid, _ := primitive.ObjectIDFromHex(pm.UserID)
	aid, _ := primitive.ObjectIDFromHex(pm.AlbumID)

	return &PhotoMetadata{
		ID:        id,
		UserID:    uid,
		AlbumID:   aid,
		BucketURL: pm.BucketURL,
		Filename:  pm.Filename,
	}
}

func toPhoto(pm *PhotoMetadata) *models.PhotoMetadata {
	return &models.PhotoMetadata{
		ID:        pm.ID.Hex(),
		UserID:    pm.UserID.Hex(),
		AlbumID:   pm.AlbumID.Hex(),
		BucketURL: pm.BucketURL,
		Filename:  pm.Filename,
	}
}
