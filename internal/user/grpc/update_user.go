package grpc

import (
	"context"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) UpdateUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	user := models.User{
		ID:       int(in.ID),
		Username: in.Username,
		Email:    in.Email,
		ImageURL: in.AvatarUrl,
	}

	userExists, err := s.service.UserExists(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, errInternal)
	}

	if userExists {
		return nil, status.Error(codes.AlreadyExists, errUsernameOrEmailIsTaken)
	}

	userData, err := s.service.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error(ctx, "update user", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	resp := userToUserPb(userData)

	return resp, nil
}
