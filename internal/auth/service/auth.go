//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mocks

package authService

import (
	"context"

	"kudago/internal/models"
)

type authService struct {
	UserDB  UserDB
}

type UserDB interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, user models.User) (bool, error)
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
	userExists, err := a.UserDB.UserExists(ctx, user)
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
