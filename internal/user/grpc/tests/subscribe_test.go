package grpc

import (
	"context"
	"testing"

	"kudago/internal/logger"
	"kudago/internal/models"
	pb "kudago/internal/user/api"
	"kudago/internal/user/grpc/tests/mocks"

	user "kudago/internal/user/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserGRPC_Subscribe(t *testing.T) {
	t.Parallel()

	subscription := models.Subscription{
		SubscriberID: 1,
		FollowsID:    2,
	}

	type expected struct {
		resp *pb.Empty
		err  error
	}

	tests := []struct {
		name      string
		req       *pb.Subscription
		setupFunc func(ctrl *gomock.Controller) *user.ServerAPI
		expected  expected
	}{
		{
			name: "success subscribe",
			req: &pb.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					Subscribe(context.Background(), subscription).
					Return(nil)
				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				resp: nil,
				err:  nil,
			},
		},
		{
			name: "not found",
			req: &pb.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					Subscribe(context.Background(), subscription).
					Return(models.ErrForeignKeyViolation)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				resp: nil,
				err:  status.Error(codes.NotFound, user.ErrUserNotFound),
			},
		},
		{
			name: "already subscribed",
			req: &pb.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					Subscribe(context.Background(), subscription).
					Return(models.ErrNothingToInsert)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				resp: nil,
				err:  status.Error(codes.AlreadyExists, user.ErrSubscriptionAlreadyExists),
			},
		},
		{
			name: "internal error",
			req: &pb.Subscription{
				SubscriberID: 1,
				FollowsID:    2,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					Subscribe(context.Background(), subscription).
					Return(models.ErrInternal)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				resp: nil,
				err:  status.Error(codes.Internal, user.ErrInternal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).Subscribe(context.Background(), tt.req)

			assert.Equal(t, tt.expected.resp, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
