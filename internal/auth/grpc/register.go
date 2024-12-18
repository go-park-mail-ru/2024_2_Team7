package auth

import (
	"context"
	"errors"

	pb "kudago/internal/auth/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.User, error) {
	user := models.User{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
		ImageURL: in.AvatarUrl,
	}

	userData, err := s.service.Register(ctx, user)
	if err != nil {
		if errors.Is(err, models.ErrEmailIsUsed) {
			return nil, status.Error(codes.AlreadyExists, ErrUsernameOrEmailIsTaken)
		}
		s.logger.Error(ctx, "register", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	resp := userToUserPb(userData)

	return resp, nil
}
