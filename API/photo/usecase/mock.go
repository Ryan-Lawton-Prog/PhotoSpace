package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
	"ryanlawton.art/photospace-api/models"
)

type PhotoUseCaseMock struct {
	mock.Mock
}

func (p PhotoUseCaseMock) UploadPhoto(ctx context.Context, user *models.User, blob *models.PhotoBlob, filename string) (string, error) {
	args := p.Called(user, blob, filename)

	return args.Get(0).(string), args.Error(1)
}

func (p PhotoUseCaseMock) FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.PhotoMetadata, models.PhotoBlob, error) {
	args := p.Called(user, photoId)

	return args.Get(0).(*models.PhotoMetadata), args.Get(1).(models.PhotoBlob), args.Error(2)
}

func (p PhotoUseCaseMock) FetchPhotoAllIDs(ctx context.Context, user *models.User) ([]string, error) {
	args := p.Called(user)

	return args.Get(0).([]string), args.Error(1)
}
