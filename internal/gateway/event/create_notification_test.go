package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kudago/internal/gateway/event/mocks"
	"kudago/internal/gateway/utils"
	"kudago/internal/logger"
	"kudago/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_CreateNotification(t *testing.T) {
	t.Parallel()

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
	}{
		{
			name: "Not found",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockNotificationServiceClient(ctrl)

				return &EventHandler{
					NotificationService: serviceMock,
					logger:              logger,
				}
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Bad request",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockNotificationServiceClient(ctrl)

				return &EventHandler{
					NotificationService: serviceMock,
					logger:              logger,
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
			tt.setupFunc(ctrl).CreateInvitationNotification(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
