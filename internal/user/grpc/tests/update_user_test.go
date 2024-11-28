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

func TestUserGRPC_UpdateUser(t *testing.T) {
	t.Parallel()

	userData := models.User{
		ID:       0,
		Username: "test",
		Email:    "test",
	}

	type expected struct {
		user *pb.User
		err  error
	}

	tests := []struct {
		name      string
		req       *pb.User
		setupFunc func(ctrl *gomock.Controller) *user.ServerAPI
		expected  expected
	}{
		{
			name: "success update",
			req: &pb.User{
				Username: userData.Username,
				Email:    userData.Email,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					UserExists(context.Background(), userData).
					Return(false, nil)

				mockUserService.EXPECT().
					UpdateUser(context.Background(), userData).
					Return(userData, nil)
				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				user: &pb.User{
					ID:       int32(userData.ID),
					Username: userData.Username,
					Email:    userData.Email,
				},
				err: nil,
			},
		},
		{
			name: "email is already taken",
			req: &pb.User{
				Username: userData.Username,
				Email:    userData.Email,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					UserExists(context.Background(), userData).
					Return(true, nil)

				return user.NewServerAPI(mockUserService, logger)
			},
			expected: expected{
				user: nil,
				err:  status.Error(codes.AlreadyExists, user.ErrUsernameOrEmailIsTaken),
			},
		},
		{
			name: "internal error",
			req: &pb.User{
				Username: userData.Username,
				Email:    userData.Email,
			},
			setupFunc: func(ctrl *gomock.Controller) *user.ServerAPI {
				mockUserService := mocks.NewMockUserService(ctrl)
				logger, _ := logger.NewLogger()

				mockUserService.EXPECT().
					UserExists(context.Background(), userData).
					Return(false, nil)

				mockUserService.EXPECT().
					UpdateUser(context.Background(), userData).
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

			actual, err := tt.setupFunc(ctrl).UpdateUser(context.Background(), tt.req)

			assert.Equal(t, tt.expected.user, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
