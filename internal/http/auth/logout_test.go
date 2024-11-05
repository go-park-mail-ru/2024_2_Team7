package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kudago/internal/http/auth/mocks"
	"kudago/internal/http/utils"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Logout(t *testing.T) {
	t.Parallel()
	randTime := time.Now().Add(time.Hour)

	tests := []struct {
		name      string
		req       *http.Request
		w         *httptest.ResponseRecorder
		setupFunc func(ctrl *gomock.Controller) *AuthHandler
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
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)
				session := models.Session{Token: "valid_token"}

				serviceMock.EXPECT().DeleteSession(gomock.Any(), session.Token).Return(nil)

				return &AuthHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Нет активной сессии",
			req:  httptest.NewRequest(http.MethodPost, "/logout", nil),
			w:    httptest.NewRecorder(),
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)

				return &AuthHandler{
					service: serviceMock,
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
			setupFunc: func(ctrl *gomock.Controller) *AuthHandler {
				serviceMock := mocks.NewMockAuthService(ctrl)
				session := models.Session{Token: "valid_token"}

				serviceMock.EXPECT().DeleteSession(gomock.Any(), session.Token).Return(models.ErrInternal)

				return &AuthHandler{
					service: serviceMock,
				}
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.setupFunc(ctrl).Logout(tt.w, tt.req)

			assert.Equal(t, tt.wantCode, tt.w.Code)

			// Additional checks can be added if needed, such as checking for cookies
		})
	}
}
