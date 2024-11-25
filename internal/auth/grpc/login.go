package auth

import (
	"context"
	"errors"

	pb "kudago/internal/auth/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.User, error) {
	creds := models.Credentials{
		Username: in.Username,
		Password: in.Password,
	}

	userData, err := s.service.CheckCredentials(ctx, creds)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errInvalidCredentials)
		}
		s.logger.Error(ctx, "login", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	user := userToUserPb(userData)

	return user, nil
}
