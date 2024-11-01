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
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
}

type SessionDB interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type ImageDB interface {
	SaveImage(ctx context.Context, header multipart.FileHeader, file multipart.File) (string, error)
	DeleteImage(ctx context.Context, imagePath string) error
}

func NewService(userDB UserDB, sessionDB SessionDB, imageDB ImageDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB, ImageDB: imageDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (models.Session, error) {
	return a.SessionDB.CheckSession(ctx, cookie)
}

func (a *authService) UpdateUser(ctx context.Context, data models.NewUserData) (models.User, error) {
	user := data.User
	oldData, err := a.UserDB.GetUserByID(ctx, user.ID)
	if err != nil {
		return models.User{}, err
	}

	if data.File != nil && data.Header != nil {
		path, err := a.ImageDB.SaveImage(ctx, *data.Header, *data.File)
		if err != nil {
			return models.User{}, err
		}
		user.ImageURL = path
	}

	user, err = a.UserDB.UpdateUser(ctx, user)
	if err != nil {
		if user.ImageURL != "" {
			a.ImageDB.DeleteImage(ctx, user.ImageURL)
		}
		return models.User{}, err
	}

	if oldData.ImageURL != "" && data.File != nil {
		err = a.ImageDB.DeleteImage(ctx, oldData.ImageURL)
	}
	return user, nil
}

func (a *authService) GetUserByID(ctx context.Context, ID int) (models.User, error) {
	return a.UserDB.GetUserByID(ctx, ID)
}

func (a *authService) CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error) {
	return a.UserDB.CheckCredentials(ctx, creds.Username, creds.Password)
}

func (a *authService) Register(ctx context.Context, data models.NewUserData) (models.User, error) {
	user := data.User

	if data.Header != nil && data.File != nil {
		path, err := a.ImageDB.SaveImage(ctx, *data.Header, *data.File)
		if err != nil {
			return models.User{}, err
		}

		user.ImageURL = path
	}
	userExists, err := a.UserDB.UserExists(ctx, user.Username, user.Email)
	if err != nil {
		return models.User{}, err
	}

	if userExists {
		return models.User{}, models.ErrEmailIsUsed
	}

	user, err = a.UserDB.AddUser(ctx, user)
	if err != nil {
		if user.ImageURL != "" {
			a.ImageDB.DeleteImage(ctx, user.ImageURL)
		}
		return models.User{}, err
	}

	return user, nil
}

func (a *authService) CreateSession(ctx context.Context, ID int) (models.Session, error) {
	return a.SessionDB.CreateSession(ctx, ID)
}

func (a *authService) DeleteSession(ctx context.Context, token string) error {
	return a.SessionDB.DeleteSession(ctx, token)
}
