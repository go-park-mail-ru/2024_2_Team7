package http

import (
	"context"
	"errors"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetUserByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.User, error) {
	userData, err := s.service.GetUserByID(ctx, int(in.ID))
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errUserNotFound)
		}
		s.logger.Error(ctx, "get user by id", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	user := userToUserPb(userData)

	return user, nil
}
