package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/user/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Parallel()

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *UserHandlers
		wantCode  int
	}{
		{
			name: "No auth",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/profile", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Bad request",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/profile", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *UserHandlers {
				serviceMock := mocks.NewMockUserServiceClient(ctrl)

				return &UserHandlers{
					UserService: serviceMock,
					logger:      logger,
				}
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).UpdateUser(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
