package events

import (
	"encoding/json"
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

func TestEventHandler_UpdateEvent(t *testing.T) {
	t.Parallel()

	logger, _ := logger.NewLogger()

	tests := []struct {
		name      string
		req       *http.Request
		setupFunc func(ctrl *gomock.Controller) *EventHandler
		wantCode  int
		wantBody  *EventResponse
	}{
		{
			name: "Успешное получение избранных событий",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events", nil)
				return req
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusForbidden,
			wantBody: &EventResponse{},
		},
		{
			name: "No id",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/events/favorites", nil)
				session := models.Session{UserID: 1, Token: "valid_token"}
				ctx := utils.SetSessionInContext(req.Context(), session)
				return req.WithContext(ctx)
			}(),
			setupFunc: func(ctrl *gomock.Controller) *EventHandler {
				serviceMock := mocks.NewMockEventServiceClient(ctrl)

				return &EventHandler{
					EventService: serviceMock,
					logger:       logger,
				}
			},
			wantCode: http.StatusBadRequest,
			wantBody: &EventResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			recorder := httptest.NewRecorder()
			tt.setupFunc(ctrl).UpdateEvent(recorder, tt.req)

			assert.Equal(t, tt.wantCode, recorder.Code)

			if tt.wantBody != nil {
				var resp EventResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &resp)
			}
		})
	}
}
