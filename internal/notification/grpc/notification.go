//go:generate mockgen -source=notification.go -destination=mocks/notification.go -package=mocks

package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	layout      = "2006-01-02 15:04:05.999999999 -0700 MST"
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
	CreateNotificationsByUserIDs(ctx context.Context, ids []int, ntf models.Notification) error
	DeleteNotification(ctx context.Context, ID int) error
	UpdateSentNotifications(ctx context.Context, IDs []int) error
}

func NewServerAPI(service NotificationService, logger *logger.Logger) *ServerAPI {
	return &ServerAPI{
		service: service,
		logger:  logger,
	}
}

func (s *ServerAPI) CreateInvitationNotification(ctx context.Context, req *pb.Notification) (*pb.Empty, error) {
	notifyAt, _ := time.Parse(layout, req.NotifyAt)

	ntf := models.Notification{
		UserID:   int(req.UserID),
		EventID:  int(req.EventID),
		NotifyAt: notifyAt,
		Message:  req.Message,
	}

	err := s.service.CreateNotification(ctx, ntf)
	if err != nil {
		s.logger.Error(ctx, "create notification", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
	return nil, nil
}

func (s *ServerAPI) CreateNotifications(ctx context.Context, req *pb.CreateNotificationsRequest) (*pb.Empty, error) {
	ids := make([]int, 0, len(req.UserIDs))
	for _, id := range req.UserIDs {
		ids = append(ids, int(id))
	}

	cleanTime := strings.Split(req.Notification.NotifyAt, " m=")[0]
	notifyAt, _ := time.Parse(layout, cleanTime)
	ntf := models.Notification{
		EventID:  int(req.Notification.EventID),
		NotifyAt: notifyAt,
		Message:  req.Notification.Message,
	}

	if err := s.service.CreateNotificationsByUserIDs(ctx, ids, ntf); err != nil {
		s.logger.Error(ctx, "create notification by user ids", err)
		return nil, status.Error(codes.Internal, errInternal)
	}
	fmt.Println(req, ntf)
	return nil, nil
}

func (s *ServerAPI) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	notifications, err := s.service.GetNotifications(ctx, int(req.UserID))
	if err != nil {
		s.logger.Error(ctx, "get notifications", err)
		return nil, status.Error(codes.Internal, errInternal)
	}

	var ids []int
	for _, n := range notifications {
		ids = append(ids, n.EventID)
	}
	err = s.service.UpdateSentNotifications(ctx, ids)
	if err != nil {
		return nil, err
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
