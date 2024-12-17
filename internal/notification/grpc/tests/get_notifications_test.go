package grpc

import (
	"context"
	"testing"
	"time"

	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/notification/api"
	notification "kudago/internal/notification/grpc"
	"kudago/internal/notification/grpc/tests/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserGRPC_GetNotifications(t *testing.T) {
	t.Parallel()

	notificationData := []models.Notification{
		{
			ID:       1,
			Message:  "test",
			EventID:  1,
			NotifyAt: time.Now(),
		},
	}

	type expected struct {
		notification *pb.GetNotificationsResponse
		err          error
	}

	tests := []struct {
		name      string
		req       *pb.GetNotificationsRequest
		setupFunc func(ctrl *gomock.Controller) *notification.ServerAPI
		expected  expected
	}{
		{
			name: "success get notification",
			req: &pb.GetNotificationsRequest{
				UserID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *notification.ServerAPI {
				mockNotificationService := mocks.NewMockNotificationService(ctrl)
				logger, _ := logger.NewLogger()

				mockNotificationService.EXPECT().
					GetNotifications(context.Background(), 1).
					Return(notificationData, nil)

				mockNotificationService.EXPECT().
					UpdateSentNotifications(context.Background(), []int{1}).
					Return(nil)
				return notification.NewServerAPI(mockNotificationService, logger)
			},
			expected: expected{
				notification: &pb.GetNotificationsResponse{
					Notifications: []*pb.Notification{
						{
							Id:       int32(notificationData[0].ID),
							Message:  notificationData[0].Message,
							EventID:  int32(notificationData[0].EventID),
							NotifyAt: notificationData[0].NotifyAt.String(),
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "internal error",
			req: &pb.GetNotificationsRequest{
				UserID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *notification.ServerAPI {
				mockNotificationService := mocks.NewMockNotificationService(ctrl)
				logger, _ := logger.NewLogger()

				mockNotificationService.EXPECT().
					GetNotifications(context.Background(), 1).
					Return(nil, models.ErrInternal)

				return notification.NewServerAPI(mockNotificationService, logger)
			},
			expected: expected{
				notification: nil,
				err:          status.Error(codes.Internal, notification.ErrInternal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).GetNotifications(context.Background(), tt.req)

			assert.Equal(t, tt.expected.notification, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
