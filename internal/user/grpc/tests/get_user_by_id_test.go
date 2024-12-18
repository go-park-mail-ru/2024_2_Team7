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

func TestUserGRPC_GetUser(t *testing.T) {
	t.Parallel()

	userData := models.User{
		ID:       1,
		Username: "test",
	}

	type expected struct {
		user *pb.User
		err  error
	}

	tests := []struct {
		name      string
		req       *pb.GetUserByIDRequest
		setupFunc func(ctrl *gomock.Controller) *user.ServerAPI
		expected  expected
	}{
		{
			name: "success get user",
			req: &pb.GetUserByIDRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(userData, nil)
				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				user: &pb.User{
					ID:       int32(userData.ID),
					Username: userData.Username,
				},
				err: nil,
			},
		},
		{
			name: "not found",
			req: &pb.GetUserByIDRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(models.User{}, models.ErrUserNotFound)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				user: nil,
				err:  status.Error(codes.NotFound, user.ErrUserNotFound),
			},
		},
		{
			name: "internal error",
			req: &pb.GetUserByIDRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(models.User{}, models.ErrInternal)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				user: nil,
				err:  status.Error(codes.Internal, user.ErrInternal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).GetUserByID(context.Background(), tt.req)

			assert.Equal(t, tt.expected.user, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
