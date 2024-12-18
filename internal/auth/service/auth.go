//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mocks

package service

import (
	"context"

	"kudago/internal/models"
)

type service struct {
	UserDB UserDB
}

type UserDB interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, user models.User) (bool, error)
}

func NewService(userDB UserDB) *service {
	return &service{UserDB: userDB}
}

func (a *service) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	return a.UserDB.GetUserByID(ctx, ID)
}

func (a *service) CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error) {
	return a.UserDB.CheckCredentials(ctx, creds.Username, creds.Password)
}

func (a *service) Register(ctx context.Context, user models.User) (models.User, error) {
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
