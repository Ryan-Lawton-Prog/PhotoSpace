package photo

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"ryanlawton.art/photospace-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var uploadPath string

type Photo struct {
	ID       primitive.ObjectID `bson:"_id",omitempty`
	UserID   primitive.ObjectID `bson:"user_id"`
	AlbumID  primitive.ObjectID `bson:"album_id"`
	PhotoURL string             `bson:"photo_url"`
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
func (pr *PhotoRepository) UploadPhoto(ctx context.Context, user *models.User, pm *models.Photo) error {

	pm.UserID = user.ID

	newPhotoId := primitive.NewObjectID()
	pm.ID = newPhotoId.Hex()

	path, err := saveImageToDisk(pm)
	if err != nil {
		log.Printf("Error saving image to disk: %s", err.Error())
		return err
	}

	model := toModel(pm)
	model.PhotoURL = path

	res, err := pr.db.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	log.Printf("Inserted photo with ID %s", res.InsertedID.(primitive.ObjectID).Hex())

	pm.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// FetchPhoto fetches a photo from the database
func (pr *PhotoRepository) FetchPhoto(ctx context.Context, user *models.User, pm *models.Photo) (*models.Photo, error) {
	pm.UserID = user.ID

	model := toModel(pm)

	res, err := pr.db.Find(ctx, bson.D{{"user_id", model.UserID}, {"_id", model.ID}})
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	photo := new(Photo)

	for res.Next(ctx) {
		err := res.Decode(photo)
		if err != nil {
			log.Printf("Error decoding photo: %s", err.Error())
			return nil, err
		}
	}

	return toPhoto(photo)
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
		photo := new(Photo)
		err := res.Decode(photo)
		if err != nil {
			log.Printf("Error decoding photo: %s", err.Error())
			return nil, err
		}

		photoIDs = append(photoIDs, photo.ID.Hex())
	}

	return photoIDs, nil
}

func toModel(p *models.Photo) *Photo {
	id, _ := primitive.ObjectIDFromHex(p.ID)
	uid, _ := primitive.ObjectIDFromHex(p.UserID)
	aid, _ := primitive.ObjectIDFromHex(p.AlbumID)

	return &Photo{
		ID:      id,
		UserID:  uid,
		AlbumID: aid,
	}
}

func toPhoto(p *Photo) (*models.Photo, error) {
	f, err := os.Open(p.PhotoURL)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return nil, err
	}
	defer f.Close()

	// Get the file size
	stat, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fileBytes := make([]byte, stat.Size())
	f.Read(fileBytes)

	return &models.Photo{
		ID:       p.ID.Hex(),
		UserID:   p.UserID.Hex(),
		AlbumID:  p.AlbumID.Hex(),
		Photo:    fileBytes,
		Filename: strings.Split(p.PhotoURL, ".")[1],
	}, nil
}

func saveImageToDisk(p *models.Photo) (string, error) {
	extension := strings.Split(p.Filename, ".")[1]
	newFileName := p.ID + "." + extension
	newPath := filepath.Join(uploadPath, newFileName)

	newFile, err := os.Create(newPath)
	if err != nil {
		return "CANT_WRITE_FILE", err
	}

	defer newFile.Close() // idempotent, okay to call twice

	if _, err := newFile.Write(p.Photo); err != nil || newFile.Close() != nil {
		return "CANT_WRITE_FILE", err
	}

	return newPath, nil
}
