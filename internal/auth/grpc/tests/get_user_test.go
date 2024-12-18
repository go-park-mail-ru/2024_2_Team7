package grpc

import (
	"context"
	"testing"

	pb "kudago/internal/auth/api"
	"kudago/internal/auth/grpc/tests/mocks"
	"kudago/internal/logger"
	"kudago/internal/models"

	auth "kudago/internal/auth/grpc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthGRPC_GetUser(t *testing.T) {
	t.Parallel()

	user := models.User{
		ID:       1,
		Username: "test",
	}

	type expected struct {
		user *pb.User
		err  error
	}

	tests := []struct {
		name      string
		req       *pb.GetUserRequest
		setupFunc func(ctrl *gomock.Controller) *auth.ServerAPI
		expected  expected
	}{
		{
			name: "success get user",
			req: &pb.GetUserRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockAuthService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(user, nil)
				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				user: &pb.User{
					ID:       int32(user.ID),
					Username: user.Username,
				},
				err: nil,
			},
		},
		{
			name: "not found",
			req: &pb.GetUserRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockAuthService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(models.User{}, models.ErrUserNotFound)

				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				user: nil,
				err:  status.Error(codes.NotFound, auth.ErrUserNotFound),
			},
		},
		{
			name: "internal error",
			req: &pb.GetUserRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockAuthService.EXPECT().
					GetUserByID(context.Background(), 1).
					Return(models.User{}, models.ErrInternal)

				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				user: nil,
				err:  status.Error(codes.Internal, auth.ErrInternal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).GetUser(context.Background(), tt.req)

			assert.Equal(t, tt.expected.user, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
