package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kudago/internal/gateway/auth/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	t.Parallel()

	type want struct {
		code     int
		username string
		email    string
		user_id  int
	}

	ctx := context.Background()
	uuid := uuid.New()
	logger, _ := logger.NewLogger()
	ctx = utils.SetRequestIDInContext(ctx, uuid.String())

	tests := []struct {
		name      string
		req       *http.Request
		w         *httptest.ResponseRecorder
		setupFunc func(ctrl *gomock.Controller) *AuthHandlers
		want      want
	}{
		{
			name: "Уже авторизован",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(`{"username": "user1", "password": "password", "email": "user1@mail.com"}`)))
				req = req.WithContext(utils.SetSessionInContext(req.Context(), models.Session{UserID: 1, Token: "abc", Expires: time.Now().Add(time.Hour)}))
				return req
			}(),
			w: httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)

				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
				}
			},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name: "Ошибка валидации JSON",
			req:  httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(`{invalid json}`))),
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)
				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
				}
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.setupFunc(ctrl).Register(tt.w, tt.req)

			assert.Equal(t, tt.want.code, tt.w.Code)
		})
	}
}
