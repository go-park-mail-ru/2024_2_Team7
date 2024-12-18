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

func TestUserGRPC_GetSubscribers(t *testing.T) {
	t.Parallel()

	usersData := []models.User{
		{
			ID:       1,
			Username: "test",
		},
	}

	type expected struct {
		resp *pb.GetSubscribersResponse
		err  error
	}

	tests := []struct {
		name      string
		req       *pb.GetSubscribersRequest
		setupFunc func(ctrl *gomock.Controller) *user.ServerAPI
		expected  expected
	}{
		{
			name: "success get subscribers",
			req: &pb.GetSubscribersRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					GetSubscribers(context.Background(), 1).
					Return(usersData, nil)
				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				resp: &pb.GetSubscribersResponse{
					Users: []*pb.User{
						{
							ID:       int32(usersData[0].ID),
							Username: usersData[0].Username,
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "internal error",
			req: &pb.GetSubscribersRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					GetSubscribers(context.Background(), 1).
					Return(nil, models.ErrInternal)

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

			actual, err := tt.setupFunc(ctrl).GetSubscribers(context.Background(), tt.req)

			assert.Equal(t, tt.expected.resp, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
