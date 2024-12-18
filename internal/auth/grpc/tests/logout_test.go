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

func TestAuthGRPC_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		req         *pb.LogoutRequest
		setupFunc   func(ctrl *gomock.Controller) *auth.ServerAPI
		expectedErr error
	}{
		{
			name: "success logout",
			req: &pb.LogoutRequest{
				Token: "valid_token",
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					DeleteSession(context.Background(), "valid_token").
					Return(nil)
				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req: &pb.LogoutRequest{
				Token: "valid_token",
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					DeleteSession(context.Background(), "valid_token").
					Return(models.ErrInternal)
				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expectedErr: status.Error(codes.Internal, auth.ErrInternal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			_, err := tt.setupFunc(ctrl).Logout(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
