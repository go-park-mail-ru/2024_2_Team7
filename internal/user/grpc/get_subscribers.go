package grpc

import (
	"context"

	"kudago/internal/models"
	pb "kudago/internal/user/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) GetSubscribers(ctx context.Context, in *pb.GetSubscribersRequest) (*pb.GetSubscribersResponse, error) {
	usersData, err := s.service.GetSubscribers(ctx, int(in.ID))
	if err != nil {
		s.logger.Error(ctx, "get subscribers", err)
		return nil, status.Error(codes.Internal, ErrInternal)
	}

	users := usersToGetSubscribersResponse(usersData)

	return users, nil
}

func usersToGetSubscribersResponse(users []models.User) *pb.GetSubscribersResponse {
	resp := &pb.GetSubscribersResponse{}

	for _, user := range users {
		userResp := userToUserPb(user)
		resp.Users = append(resp.Users, userResp)
	}
	return resp
}
