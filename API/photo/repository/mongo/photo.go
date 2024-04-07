package photo

import (
	"context"
	"image"
	"log"

	"ryanlawton.art/photospace-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Photo struct {
	ID      primitive.ObjectID `bson:"_id",omitempty`
	UserID  primitive.ObjectID `bson:"user_id"`
	AlbumID primitive.ObjectID `bson:"album_id"`
	Photo   image.Image        `bson:"photo"`
	Title   string             `bson:"title"`
}

type PhotoRepository struct {
	db *mongo.Collection
}

func NewPhotoRepository(db *mongo.Database, collection string) *PhotoRepository {
	return &PhotoRepository{
		db: db.Collection(collection),
	}
}

func (pr *PhotoRepository) UploadPhoto(ctx context.Context, user *models.User, pm *models.Photo) error {
	pm.UserID = user.ID

	model := toModel(pm)

	res, err := pr.db.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	log.Println(res.InsertedID)

	pm.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (pr *PhotoRepository) FetchPhoto(ctx context.Context, user *models.User, pm *models.Photo) (*models.Photo, error) {
	pm.UserID = user.ID

	model := toModel(pm)

	log.Println(model)

	res, err := pr.db.Find(ctx, model)
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	log.Println(res)

	photo := new(Photo)

	for res.Next(ctx) {
		err := res.Decode(photo)
		if err != nil {
			return nil, err
		}
	}

	log.Println(photo)

	return toPhoto(photo), nil
}

func toModel(p *models.Photo) *Photo {
	uid, _ := primitive.ObjectIDFromHex(p.UserID)

	return &Photo{
		UserID:  uid,
		AlbumID: uid,
		Photo:   p.Photo,
		Title:   p.Title,
	}
}

func toPhoto(p *Photo) *models.Photo {
	return &models.Photo{
		ID:      p.ID.Hex(),
		UserID:  p.UserID.Hex(),
		AlbumID: p.AlbumID.Hex(),
		Photo:   p.Photo,
		Title:   p.Title,
	}
}

func toPhotos(ps []*Photo) []*models.Photo {
	out := make([]*models.Photo, len(ps))

	for i, b := range ps {
		out[i] = toPhoto(b)
	}

	return out
}
