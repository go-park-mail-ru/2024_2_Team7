package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pb "kudago/internal/auth/api"
	auth "kudago/internal/auth/grpc"
	"kudago/internal/gateway/auth/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_Logout(t *testing.T) {
	t.Parallel()
	randTime := time.Now().Add(time.Hour)
	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		w         *httptest.ResponseRecorder
		setupFunc func(ctrl *gomock.Controller) *AuthHandlers
		wantCode  int
	}{
		{
			name: "Успешный выход из системы",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer([]byte(`{"username": "user1", "password": "password"}`)))
				req = req.WithContext(utils.SetSessionInContext(req.Context(), models.Session{UserID: 1, Token: "valid_token", Expires: randTime}))
				return req
			}(),
			w: httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)
				session := &pb.Session{Token: "valid_token"}

				serviceMock.EXPECT().Logout(gomock.Any(), &pb.LogoutRequest{Token: session.Token}).Return(nil, nil)

				return &AuthHandlers{
					AuthService: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Нет активной сессии",
			req:  httptest.NewRequest(http.MethodPost, "/logout", nil),
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)

				return &AuthHandlers{
					AuthService: serviceMock,
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Ошибка при удалении сессии",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/logout", nil)
				session := models.Session{Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			w: httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)
				session := &pb.Session{Token: "valid_token"}

				serviceMock.EXPECT().Logout(gomock.Any(), &pb.LogoutRequest{Token: session.Token}).Return(nil, status.Error(codes.NotFound, auth.ErrInternal))

				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.setupFunc(ctrl).Logout(tt.w, tt.req)

			assert.Equal(t, tt.wantCode, tt.w.Code)
		})
	}
}
