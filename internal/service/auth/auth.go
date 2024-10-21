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
	UserExists(ctx context.Context, username string) bool
	AddUser(ctx context.Context, user *models.User) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) bool
}

type SessionDB interface {
	CheckSession(ctx context.Context, cookie string) (*models.Session, bool)
	CreateSession(ctx context.Context, ID int) *models.Session
	DeleteSession(ctx context.Context, token string)
}

func NewService(userDB UserDB, sessionDB SessionDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (*models.Session, bool) {
	return a.SessionDB.CheckSession(ctx, cookie)
}

func (a *authService) UserExists(ctx context.Context, username string) bool {
	return a.UserDB.UserExists(ctx, username)
}

func (a *authService) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	return a.UserDB.AddUser(ctx, user)
}

func (a *authService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return a.UserDB.GetUserByUsername(ctx, username)
}

func (a *authService) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	return a.UserDB.GetUserByID(ctx, ID)
}

func (a *authService) CheckCredentials(ctx context.Context, creds models.Credentials) bool {
	return a.UserDB.CheckCredentials(ctx, creds.Username, creds.Password)
}

func (a *authService) Register(ctx context.Context, user models.User) (models.User, error) {
	if a.UserDB.UserExists(ctx, user.Username) {
		return models.User{}, models.ErrUsernameIsUsed
	}

	user, err := a.UserDB.AddUser(ctx, &user)
	if err != nil {
		return models.User{}, models.ErrEmailIsUsed
	}

	return user, nil
}

func (a *authService) CreateSession(ctx context.Context, ID int) *models.Session {
	return a.SessionDB.CreateSession(ctx, ID)
}

func (a *authService) DeleteSession(ctx context.Context, token string) {
	a.SessionDB.DeleteSession(ctx, token)
}
