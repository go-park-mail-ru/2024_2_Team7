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

func TestAuthGRPC_CheckSession(t *testing.T) {
	t.Parallel()

	session := models.Session{
		UserID: 1,
		Token:  "test",
	}

	type expected struct {
		session *pb.Session
		err     error
	}

	tests := []struct {
		name      string
		req       *pb.CheckSessionRequest
		setupFunc func(ctrl *gomock.Controller) *auth.ServerAPI
		expected  expected
	}{
		{
			name: "success check session",
			req: &pb.CheckSessionRequest{
				Cookie: "cookie",
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					CheckSession(context.Background(), "cookie").
					Return(session, nil)
				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				session: &pb.Session{
					UserID:  int32(session.UserID),
					Token:   session.Token,
					Expires: "0001-01-01T00:00:00Z",
				},
				err: nil,
			},
		},
		{
			name: "not found",
			req: &pb.CheckSessionRequest{
				Cookie: "cookie",
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					CheckSession(context.Background(), "cookie").
					Return(models.Session{}, models.ErrUserNotFound)

				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				session: nil,
				err:     status.Error(codes.NotFound, auth.ErrUserNotFound),
			},
		},
		{
			name: "internal error",
			req: &pb.CheckSessionRequest{
				Cookie: "cookie",
			},
			setupFunc: func(ctrl *gomock.Controller) *auth.ServerAPI {
				mockSessionManager := mocks.NewMockSessionManager(ctrl)
				mockAuthService := mocks.NewMockAuthService(ctrl)
				logger, _ := logger.NewLogger()

				mockSessionManager.EXPECT().
					CheckSession(context.Background(), "cookie").
					Return(models.Session{}, models.ErrInternal)

				return auth.NewServerAPI(mockAuthService, mockSessionManager, logger)
			},
			expected: expected{
				session: nil,
				err:     status.Error(codes.Internal, auth.ErrInternal),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			actual, err := tt.setupFunc(ctrl).CheckSession(context.Background(), tt.req)

			assert.Equal(t, tt.expected.session, actual)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
