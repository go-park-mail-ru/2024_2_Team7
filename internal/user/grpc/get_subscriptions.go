package grpc

import (
	"context"
	"errors"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetSubscriptions(ctx context.Context, in *pb.GetSubscriptionsRequest) (*pb.GetSubscriptionsResponse, error) {
	usersData, err := s.service.GetSubscriptions(ctx, int(in.ID))
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound)
		}
		s.logger.Error(ctx, "get subscriptions", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	users := usersToUsersPb(usersData)

	return users, nil
}

func usersToUsersPb(users []models.User) *pb.GetSubscriptionsResponse {
	resp := &pb.GetSubscriptionsResponse{}

	for _, user := range users {
		userResp := userToUserPb(user)
		resp.Users = append(resp.Users, userResp)
	}
	return resp
}
