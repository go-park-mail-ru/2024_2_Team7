package http

import (
	"context"
	"errors"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) Unsubscribe(ctx context.Context, in *pb.Subscription) (*pb.Empty, error) {
	subscription := subscriptionPBToSubscription(in)

	err := s.service.Unsubscribe(ctx, subscription)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, errUserNotFound)
		}
		s.logger.Error(ctx, "get user by id", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}
