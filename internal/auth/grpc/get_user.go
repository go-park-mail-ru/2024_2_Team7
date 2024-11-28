package auth

import (
	"context"
	"errors"

	pb "kudago/internal/auth/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	userData, err := s.service.GetUserByID(ctx, int(in.ID))
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound)
		}
		s.logger.Error(ctx, "get user", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	user := userToUserPb(userData)

	return user, nil
}
