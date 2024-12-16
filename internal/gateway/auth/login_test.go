package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	randTime := time.Now().Add(time.Hour).Format(time.RFC3339)

	creds := &pb.LoginRequest{
		Username: "user1",
		Password: "password",
	}

	body, err := json.Marshal(creds)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/login", bytes.NewBuffer(body))
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *http.Request
		w         *httptest.ResponseRecorder
		setupFunc func(ctrl *gomock.Controller) *AuthHandlers
		want      want
	}{
		{
			name: "Успешный логин",
			req:  req,
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)

				serviceMock.EXPECT().Login(gomock.Any(), creds).
					Return(&pb.User{
						ID:       1,
						Username: "user1",
						Email:    "user1@mail.ru",
					}, nil)

				serviceMock.EXPECT().CreateSession(gomock.Any(), &pb.CreateSessionRequest{ID: 1}).
					Return(&pb.Session{
						UserID:  1,
						Token:   "token",
						Expires: randTime,
					}, nil)

				return &AuthHandlers{
					AuthService: serviceMock,
					logger:      logger,
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
			setupFunc: func(ctrl *gomock.Controller) *AuthHandlers {
				serviceMock := mocks.NewMockAuthServiceClient(ctrl)

				serviceMock.EXPECT().Login(gomock.Any(), creds).
					Return(&pb.User{}, status.Error(codes.NotFound, auth.ErrUserNotFound))

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
			name: "Уже авторизован",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{"username": "user1", "password": "password"}`)))
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
			req:  httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{invalid json}`))),
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
