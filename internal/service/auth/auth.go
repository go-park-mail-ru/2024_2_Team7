package authService

import (
	"context"

	"kudago/internal/models"
)

type authService struct {
	UserDB    UserDB
	SessionDB SessionDB
	CsrfDB    CsrfDB
}

type UserDB interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	CheckCredentials(ctx context.Context, username string, password string) (models.User, error)
	UserExists(ctx context.Context, username, email string) (bool, error)
}

type SessionDB interface {
	CheckSession(ctx context.Context, cookie string) (models.Session, bool)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
	DeleteSession(ctx context.Context, token string)
}

type CsrfDB interface {
	CreateCSRF(ctx context.Context, encryptionKey []byte, s *models.Session) (string, error)
	CheckCSRF(ctx context.Context, encryptionKey []byte, s *models.Session, inputToken string) (bool, error)
}

func NewService(userDB UserDB, sessionDB SessionDB, csrfDB CsrfDB) authService {
	return authService{UserDB: userDB, SessionDB: sessionDB, CsrfDB: csrfDB}
}

func (a *authService) CheckSession(ctx context.Context, cookie string) (models.Session, bool) {
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

func (a *authService) DeleteSession(ctx context.Context, token string) {
	a.SessionDB.DeleteSession(ctx, token)
}
func (a *authService) CreateCSRF(ctx context.Context, encryptionKey []byte, s *models.Session) (string, error) {
	return a.CsrfDB.CreateCSRF(ctx, encryptionKey, s)
}

func (a *authService) CheckCSRF(ctx context.Context, encryptionKey []byte, s *models.Session, inputToken string) (bool, error) {
	return a.CsrfDB.CheckCSRF(ctx, encryptionKey, s, inputToken)
}
