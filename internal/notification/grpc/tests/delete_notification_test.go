package grpc

import (
	"context"
	"testing"

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

func TestUserGRPC_DeleteNotification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *pb.DeleteNotificationRequest
		setupFunc func(ctrl *gomock.Controller) *notification.ServerAPI
		err       error
	}{
		{
			name: "success delete notification",
			req: &pb.DeleteNotificationRequest{
				Id: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *notification.ServerAPI {
				mockNotificationService := mocks.NewMockNotificationService(ctrl)
				logger, _ := logger.NewLogger()

				mockNotificationService.EXPECT().
					DeleteNotification(context.Background(), 1).
					Return(nil)

				return notification.NewServerAPI(mockNotificationService, logger)
			},
			err: nil,
		},
		{
			name: "internal error",
			req: &pb.DeleteNotificationRequest{
				Id: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *notification.ServerAPI {
				mockNotificationService := mocks.NewMockNotificationService(ctrl)
				logger, _ := logger.NewLogger()

				mockNotificationService.EXPECT().
					DeleteNotification(context.Background(), 1).
					Return(models.ErrInternal)

				return notification.NewServerAPI(mockNotificationService, logger)
			},
			err: status.Error(codes.Internal, notification.ErrInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, err := tt.setupFunc(ctrl).DeleteNotification(context.Background(), tt.req)

			assert.Equal(t, tt.err, err)
		})
	}
}
