package auth

import (
	"context"
	"errors"

	pb "kudago/internal/auth/api"
	"kudago/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) CheckSession(ctx context.Context, in *pb.CheckSessionRequest) (*pb.User, error) {
	userData, err := s.service.GetUserByID(ctx, int(in.ID))
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errUserNotFound)
		}
		s.logger.Error(ctx, "check session", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	user := userToUserPb(userData)

	return user, nil
}
