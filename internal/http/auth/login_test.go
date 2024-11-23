package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kudago/internal/http/auth/mocks"
	"kudago/internal/http/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login(t *testing.T) {
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
	randTime := time.Now().Add(time.Hour)

	creds := models.Credentials{
		Username: "user1",
		Password: "password",
	}

	body, err := json.Marshal(creds)
	assert.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/login", bytes.NewBuffer(body))
	assert.NoError(t, err)

	tests := []struct {
		name      string
		req       *http.Request
		w         *httptest.ResponseRecorder
		setupFunc func(ctrl *gomock.Controller) *AuthHandler
		want      want
	}{
		{
			name: "Успешный логин",
			req:  req,
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)

				serviceMock.EXPECT().CheckCredentials(gomock.Any(), creds).
					Return(models.User{
						ID:       1,
						Username: "user1",
						Email:    "user1@mail.ru",
					}, nil)
				serviceMock.EXPECT().CreateSession(gomock.Any(), 1).
					Return(models.Session{
						UserID:  1,
						Token:   "token",
						Expires: randTime,
					}, nil)

				return &AuthHandler{
					service: serviceMock,
					logger:  logger,
				}
			},
			want: want{
				code:     http.StatusOK,
				user_id:  1,
				username: "user1",
				email:    "user1@mail.ru",
			},
		},
		{
			name: "Неверные данные",
			req:  httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{"username": "user1", "password": "password"}`))),
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)

				serviceMock.EXPECT().CheckCredentials(gomock.Any(), creds).
					Return(models.User{}, models.ErrUserNotFound)

				return &AuthHandler{
					service: serviceMock,
					logger:  logger,
				}
			},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name: "Уже авторизован",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{"username": "user1", "password": "password"}`)))
				req = req.WithContext(utils.SetSessionInContext(req.Context(), models.Session{UserID: 1, Token: "abc", Expires: randTime}))
				return req
			}(),
			w: httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)

				return &AuthHandler{
					service: serviceMock,
					logger:  logger,
				}
			},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name: "Ошибка валидации JSON",
			req:  httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{invalid json}`))),
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)
				return &AuthHandler{
					service: serviceMock,
					logger:  logger,
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

			tt.setupFunc(ctrl).Login(tt.w, tt.req)

			assert.Equal(t, tt.want.code, tt.w.Code)
		})
	}
}

// var resp UserResponse
// err := json.Unmarshal(tt.w.Body., &resp)
// assert.NoError(t, err)

// assert.Equal(t, tt.want.email, resp.Email)
// assert.Equal(t, tt.want.username, resp.Username)
// assert.Equal(t, tt.want.user_id, resp.ID)
