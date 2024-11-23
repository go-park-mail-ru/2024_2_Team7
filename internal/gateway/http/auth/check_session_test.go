package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/http/auth/mocks"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_CheckSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *AuthHandler
		wantCode  int
		wantBody  *AuthResponse
	}{
		{
			name: "Успешная проверка сессии",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/session", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)
				user := models.User{
					ID:       1,
					Username: "user1",
					Email:    "user1@mail.ru",
				}

				serviceMock.EXPECT().GetUserByID(gomock.Any(), 1).Return(user, nil)

				return &AuthHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusOK,
			wantBody: &AuthResponse{
				User: UserResponse{
					ID:       1,
					Username: "user1",
					Email:    "user1@mail.ru",
				},
			},
		},
		{
			name: "Нет активной сессии",
			req:  httptest.NewRequest(http.MethodGet, "/session", nil),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)
				return &AuthHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Пользователь не найден",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/session", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)

				serviceMock.EXPECT().GetUserByID(gomock.Any(), 1).Return(models.User{}, models.ErrUserNotFound)

				return &AuthHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).CheckSession(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp AuthResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
