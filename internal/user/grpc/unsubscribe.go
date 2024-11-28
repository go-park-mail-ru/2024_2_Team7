package grpc

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
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound)
		}
		s.logger.Error(ctx, "unsubscribe", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	return nil, nil
}
