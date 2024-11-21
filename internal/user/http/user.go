package http

import (
	"context"

	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/user/api"
)

type ServerAPI struct {
	pb.UnimplementedUserServiceServer
	service UserService
	logger  *logger.Logger
}

type UserService interface {
	GetUserByID(ctx context.Context, ID int) (models.User, error)
	Subscribe(ctx context.Context, subscription models.Subscription) error
	Unsubscribe(ctx context.Context, subscription models.Subscription) error
	GetSubscriptions(ctx context.Context, ID int) ([]models.User, error)
}

func NewServerAPI(service UserService, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		logger:  logger,
	}
}

func userToUserPb(userData models.User) *pb.User {
	return &pb.User{
		ID:        int32(userData.ID),
		Username:  userData.Username,
		Email:     userData.Email,
		AvatarUrl: userData.ImageURL,
	}
}

func subscriptionPBToSubscription(subscription *pb.Subscription) models.Subscription {
	return models.Subscription{
		SubscriberID: int(subscription.SubscriberID),
		FollowsID:    int(subscription.FollowsID),
	}
}
