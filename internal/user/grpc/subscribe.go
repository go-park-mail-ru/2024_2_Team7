package grpc

import (
	"context"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) Subscribe(ctx context.Context, in *pb.Subscription) (*pb.Empty, error) {
	subscription := subscriptionPBToSubscription(in)

	err := s.service.Subscribe(ctx, subscription)
	if err != nil {
		s.logger.Error(ctx, "subscribe", err)
		switch err {
		case models.ErrForeignKeyViolation:
			return nil, status.Error(codes.NotFound, ErrUserNotFound)
		case models.ErrNothingToInsert:
			return nil, status.Error(codes.AlreadyExists, ErrSubscriptionAlreadyExists)
		default:
			return nil, status.Error(codes.Internal, ErrInternal)
		}
	}

	return nil, nil
}
