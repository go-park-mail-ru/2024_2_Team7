package authService

import (
	"context"
	"mime/multipart"

	"kudago/internal/models"
)

type authService struct {
	UserDB    UserDB
	SessionDB SessionDB
	ImageDB   ImageDB
}

type UserDB interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
	UpdateUser(ctx context.Context, user models.User) error
}

type SessionDB interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type ImageDB interface {
	SaveImage(ctx context.Context, header multipart.FileHeader, file multipart.File) (string, error)
}

func NewService(userDB UserDB, sessionDB SessionDB, imageDB ImageDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB, ImageDB: imageDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	return a.SessionDB.CheckSession(ctx, cookie)
}

func (a *authService) UpdateUser(ctx context.Context, user models.User) error {
	return a.UserDB.UpdateUser(ctx, user)
}

func (a *authService) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	return a.UserDB.GetUserByID(ctx, ID)
}

func (a *authService) CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error) {
	return a.UserDB.CheckCredentials(ctx, creds.Username, creds.Password)
}

func (a *authService) Register(ctx context.Context, registerDTO models.RegisterDTO) (models.User, error) {
	path, err := a.ImageDB.SaveImage(ctx, registerDTO.Header, registerDTO.File)
	if err != nil {
		return models.User{}, err
	}

	user := registerDTO.User
	user.ImageURL = path
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
