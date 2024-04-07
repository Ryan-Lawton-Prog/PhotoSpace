package auth

import (
	"context"

	"ryanlawton.art/photospace-api/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, username, password string) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
