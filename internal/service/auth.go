package service

import (
	"context"
	"kudago/internal/models"
	"kudago/internal/repository"
)

type authService struct {
	UserDB    iUserDB
	SessionDB iSessionDB
}

func NewAuthService(userDB *repository.UserDB, sessionDB *repository.SessionDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (*models.Session, bool) {
	return a.SessionDB.CheckSession(ctx, cookie)
}

func (a *authService) UserExists(ctx context.Context,username string) bool {
	return a.UserDB.UserExists(ctx,username)
}

func (a *authService) AddUser(ctx context.Context,user *models.User) error {
	return a.UserDB.AddUser(ctx,user)
}

func (a *authService) GetUser(ctx context.Context,username string) models.User {
	user := a.UserDB.GetUser(ctx,username)
	user.Password = ""
	return user
}

func (a *authService) CheckCredentials(ctx context.Context,creds models.Credentials) bool {
	return a.UserDB.CheckCredentials(ctx,creds.Username, creds.Password)
}
func (a *authService) Register(ctx context.Context,user models.User) error {
	if a.UserDB.UserExists(ctx,user.Username) {
		return models.ErrUserAlreadyExists
	}

	if err := a.UserDB.AddUser(ctx,&user); err != nil {
		return models.ErrEmailIsUsed
	}

	return nil
}

func (a *authService) CreateSession(ctx context.Context,username string) *models.Session {
	return a.SessionDB.CreateSession(ctx,username)
}

func (a *authService) DeleteSession(ctx context.Context,username string) {
	a.SessionDB.DeleteSession(ctx,username)
}
