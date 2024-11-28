//go:generate mockgen -source=auth.go -destination=mocks/auth.go -package=mocks

package auth

import (
	"context"

	pb "kudago/internal/auth/api"
	"kudago/internal/logger"
	"kudago/internal/models"
)

type ServerAPI struct {
	pb.UnimplementedAuthServiceServer
	service        AuthService
	sessionManager SessionManager
	logger         *logger.Logger
}

type AuthService interface {
	CheckCredentials(ctx context.Context, creds models.Credentials) (models.User, error)
	Register(ctx context.Context, user models.User) (models.User, error)
	GetUserByID(ctx context.Context, ID int) (models.User, error)
}

type SessionManager interface {
	DeleteSession(ctx context.Context, token string) error
	CheckSession(ctx context.Context, cookie string) (models.Session, error)
	CreateSession(ctx context.Context, ID int) (models.Session, error)
}

func NewServerAPI(service AuthService, sessionManager SessionManager, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		sessionManager: sessionManager,
		service:        service,
		logger:         logger,
	}
}

func userToUserPb(userData models.User) *pb.User {
	return &pb.User{
		ID:        int32(userData.ID),
		Username:  userData.Username,
		Email:     userData.Email,
		AvatarUrl: userData.ImageURL,
	}
}
