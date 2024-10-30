package authService

import (
	"context"

	"kudago/internal/models"
)

type authService struct {
	UserDB    UserDB
	SessionDB SessionDB
}

type UserDB interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
}

type SessionDB interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

func NewService(userDB UserDB, sessionDB SessionDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	return a.SessionDB.CheckSession(ctx, cookie)
}

func (a *authService) AddUser(ctx context.Context, user models.User) (models.User, error) {
	return a.UserDB.AddUser(ctx, user)
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
	return a.UserDB.AddUser(ctx, user)
}

func (a *authService) CreateSession(ctx context.Context, ID int) (models.Session, error) {
	return a.SessionDB.CreateSession(ctx, ID)
}

func (a *authService) DeleteSession(ctx context.Context, token string) error {
	return a.SessionDB.DeleteSession(ctx, token)
}
