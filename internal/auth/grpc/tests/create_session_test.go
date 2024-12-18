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

func TestAuthGRPC_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	session := models.Session{
		UserID: 1,
		Token:  "test",
	}

	tests := []struct {
		name        string
		req         *pb.CreateSessionRequest
		setupFunc   func(ctrl *gomock.Controller) *auth.ServerAPI
		expectedErr error
	}{
		{
			name: "success create session",
			req: &pb.CreateSessionRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					CreateSession(context.Background(), 1).
					Return(session, nil)
				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expectedErr: nil,
		},
		{
			name: "internal error",
			req: &pb.CreateSessionRequest{
				ID: 1,
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					CreateSession(context.Background(), 1).
					Return(models.Session{}, models.ErrInternal)
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

			_, err := tt.setupFunc(ctrl).CreateSession(context.Background(), tt.req)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
