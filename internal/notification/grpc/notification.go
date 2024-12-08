//go:generate mockgen -source=csat.go -destination=mocks/csat.go -package=mocks

package grpc

import (
	"context"
	"time"

	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errInternal = "internal error"
)

type ServerAPI struct {
	pb.UnimplementedNotificationServiceServer
	service NotificationService
	logger  *logger.Logger
}

type NotificationService interface {
	GetNotifications(ctx context.Context, userID int) ([]models.Notification, error)
	CreateNotification(ctx context.Context, notification models.Notification) error
	DeleteNotification(ctx context.Context, ID int) error
}

func NewServerAPI(service NotificationService, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		logger:  logger,
	}
}

func (s *ServerAPI) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.Empty, error) {
	notification := toNotificationModel(req.Notification)
	err := s.service.CreateNotification(ctx, notification)
	if err != nil {
		s.logger.Error(ctx, "create notification", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
	return nil, nil
}

func (s *ServerAPI) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	notifications, err := s.service.GetNotifications(ctx, int(req.UserID))
	if err != nil {
		s.logger.Error(ctx, "get notifications", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	resp := toGetNotificationsResponse(notifications)
	return resp, nil
}

func toGetNotificationsResponse(notifications []models.Notification) *pb.GetNotificationsResponse {
	notificationsPB := make([]*pb.Notification, 0, len(notifications))

	for _, ntf := range notifications {
		temp := &pb.Notification{
			Id:       int32(ntf.ID),
			UserID:   int32(ntf.UserID),
			EventID:  int32(ntf.EventID),
			Message:  ntf.Message,
			NotifyAt: ntf.NotifyAt.String(),
		}
		notificationsPB = append(notificationsPB, temp)
	}

	return &pb.GetNotificationsResponse{
		Notifications: notificationsPB,
	}
}

func (s *ServerAPI) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.Empty, error) {
	err := s.service.DeleteNotification(ctx, int(req.Id))
	if err != nil {
		s.logger.Error(ctx, "delete notification", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	return nil, nil
}

func toNotificationModel(notification *pb.Notification) models.Notification {
	notifyAt, _ := time.Parse(time.RFC3339, notification.NotifyAt)

	return models.Notification{
		ID:       int(notification.Id),
		UserID:   int(notification.UserID),
		EventID:  int(notification.EventID),
		NotifyAt: notifyAt,
		Message:  notification.Message,
	}
}
