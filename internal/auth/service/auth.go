//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mocks

package authService

import (
	"context"

	"kudago/internal/models"
)

type authService struct {
	UserDB  UserDB
	ImageDB ImageDB
}

type UserDB interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	CheckUsername(ctx context.Context, username string, ID int) (bool, error)
	CheckEmail(ctx context.Context, email string, ID int) (bool, error)
}

type ImageDB interface {
	SaveImage(ctx context.Context, media models.MediaFile) (string, error)
	DeleteImage(ctx context.Context, imagePath string) error
}

func NewService(userDB UserDB) authService {
	return authService{UserDB: userDB}
}

func (a *authService) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	return a.UserDB.GetUserByID(ctx, ID)
}

func (a *authService) CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error) {
	return a.UserDB.CheckCredentials(ctx, creds.Username, creds.Password)
}

func (a *authService) Register(ctx context.Context, user models.User) (models.User, error) {
	userExists, err := a.UserDB.UserExists(ctx, user.Username, user.Email)
	if err != nil {
		return models.User{}, err
	}

	if userExists {
		return models.User{}, models.ErrEmailIsUsed
	}

	user, err = a.UserDB.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
