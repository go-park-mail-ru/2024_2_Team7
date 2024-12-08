package handlers

import (
	"kudago/internal/logger"
	pb "kudago/internal/notification/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationHandler struct {
	NotificationService pb.NotificationServiceClient
	logger              *logger.Logger
}

func NewHandlers(notificationServiceAddr string, logger *logger.Logger) (*NotificationHandler, error) {
	notificationConn, err := grpc.NewClient(notificationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &NotificationHandler{
		NotificationService: pb.NewNotificationServiceClient(notificationConn),
		logger:              logger,
	}, nil
}
