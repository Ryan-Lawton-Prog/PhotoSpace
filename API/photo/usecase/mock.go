package usecase

import (
	"context"
	"image"

	"github.com/stretchr/testify/mock"
	"ryanlawton.art/photospace-api/models"
)

type PhotoUseCaseMock struct {
	mock.Mock
}

func (p PhotoUseCaseMock) UploadPhoto(ctx context.Context, user *models.User, photo image.Image, extension string) error {
	args := p.Called(user, photo, extension)

	return args.Error(0)
}

func (p PhotoUseCaseMock) FetchPhoto(ctx context.Context, user *models.User, photoId string) (*models.Photo, error) {
	args := p.Called(user, photoId)

	return args.Get(0).(*models.Photo), args.Error(1)
}
