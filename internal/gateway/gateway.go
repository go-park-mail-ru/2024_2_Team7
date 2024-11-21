package gateway

import (
	auth "kudago/internal/auth/api"
	pbAuth "kudago/internal/auth/api"
	pbUser "kudago/internal/user/api"
	user "kudago/internal/user/api"

	"kudago/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	AuthService pbAuth.AuthServiceClient
	UserService pbUser.UserServiceClient
	Logger      *logger.Logger
}

func NewGateway(authServiceAddr string, userServiceAddr string, logger *logger.Logger) (*Gateway, error) {
	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	userConn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Gateway{
		AuthService: auth.NewAuthServiceClient(authConn),
		UserService: user.NewUserServiceClient(userConn),
		Logger:      logger,
	}, nil
}
